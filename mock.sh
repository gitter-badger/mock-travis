#!/bin/bash

function add_extra_repo()
{
        if [ -n "${yml_mock_travis_packages_extra_repo}" ]
        then {
                color_printlin cyan "Adding ${yml_mock_travis_packages_extra_repo} as extra repository"
                sed -i '$ d' /etc/mock/${yml_mock_travis_mock_config}.cfg
                add_mock_config "[extra-local]"
                add_mock_config "name=extra-local"
                add_mock_config "baseurl=${yml_mock_travis_packages_extra_repo}"
                add_mock_config "gpgcheck=0"
                add_mock_config "\"\"\""
                check_status "Adding extra repository succeeded" "Adding extra repository failed"
        } else {
                color_printlin cyan "Extra repository is not set"
                color_printlin green "No extra repository will be used on building packages"
        }
        fi
}

function add_local_repo()
{
        sed -i '$ d' /etc/mock/${yml_mock_travis_mock_config}.cfg
        add_mock_config "[mock-local]"
        add_mock_config "name=mock-local"
        add_mock_config "baseurl=file:///home/"
        add_mock_config "gpgcheck=0"
        add_mock_config "\"\"\""
}

function add_mock_config()
{
        echo $1 >> /etc/mock/${yml_mock_travis_mock_config}.cfg
}

function build_buildrequires()
{
        if [ -n "${yml_mock_travis_packages_buildrequires}" ]
        then {
                clean_dnf_cache
                build_for_local ${yml_mock_travis_packages_buildrequires}
        } else {
                clean_dnf_cache
                build_for_git
        }
        fi
}

function build_for_git()
{
        GIT_URL=https://github.com/${yml_mock_travis_packages_buildrequires_git}
        color_printlin cyan "Start setting git repository"
        dnf -y install git > /dev/null 2>&1
        cd /home/ > /dev/null 2>&1
        git clone ${GIT_URL} /home/GIT > /dev/null 2>&1
        check_status "Setting git repository succeeded" "Setting git repository failed"
        mock_build_chain `dirname \`find /home/GIT -name "*.spec" | grep GIT \` | tr '\n' ' '`
}

function build_for_local()
{
        mock_build_chain `find_dir $1 | tr '\n' ' '`
}

function build_target_pkg()
{
        clean_dnf_cache
        build_for_local ${yml_mock_travis_packages_name}
}

function check_status()
{
	if [ $? != 0 ]
	then {
	        echo -e "\e[31m\e[1m$2.\e[0m"
		exit 1
	} else {
	        echo -e "\e[32m\e[1m$1.\e[0m"
	}
	fi
}

function clean_dnf_cache()
{
        color_printlin cyan "Start cleaning dnf package manager cache"
        /usr/bin/mock -r ${yml_mock_travis_mock_config} --dnf-cmd clean all > /dev/null 2>&1
        check_status "Clean dnf package manager cache succeeded" "Clean dnf package manager cache failed"
}

function color_printlin()
{
        case "$1" in
                "cyan")
                        echo -e "\n\e[36m\e[1m$2...\e[0m"
                ;;
                "green")
                        echo -e "\n\e[32m\e[1m$2.\e[0m"
                ;;
                "red")
                        echo -e "\n\e[31m\e[1m$2.\e[0m"
                ;;
        esac
}

function find_dir()
{
        if [ -n "${yml_mock_travis_packages_directory}" ]
        then {
                cd /home/
                cd ${yml_mock_travis_packages_directory}/
                LOCATION=`pwd`
        } else {
                LOCATION="/home"
        }
        fi
        
        for PKG_NAME in $@; {
                dirname `find ${LOCATION} -name "${PKG_NAME}.spec"`
        }
}

function mock_build()
{
        for BUILD_LIST in $@; {
                color_printlin cyan "Start downloading $1 source files"
                spectool -g $1.spec > /dev/null 2>&1
                check_status "Download source succeeded" "Download source failed"
                color_printlin cyan "Start building $1 SRPM"
                /usr/bin/mock -r ${yml_mock_travis_mock_config} --resultdir ./ --buildsrpm --sources ./ --spec $1.spec > /dev/null
                check_status "Build $1 SRPM succeeded" "Build $1 SRPM failed"
                color_printlin cyan "Start building $1 binary RPM"
                /usr/bin/mock -r ${yml_mock_travis_mock_config} --resultdir ./ --rebuild `ls | grep *src.rpm` > /dev/null
                check_status "Build $1 succeeded" "Build $1 failed"
        }
}

function mock_build_chain()
{
        for PKG_DIR in $@; {
                cd ${PKG_DIR}/
                PKG_LIST=`basename -s .spec -a \`find ${PKG_DIR} -name "*.spec"\` | tr '\n' ' '`
                mock_build ${PKG_LIST}
        }
}

function read_config()
{
        local prefix=$2
        local s='[[:space:]]*' w='[a-zA-Z0-9_]*' fs=$(echo @|tr @ '\034')
        sed -ne "s|^\($s\)\($w\)$s:$s\"\(.*\)\"$s\$|\1$fs\2$fs\3|p" \
                -e "s|^\($s\)\($w\)$s:$s\(.*\)$s\$|\1$fs\2$fs\3|p"  $1 |
        awk -F$fs '{
                indent = length($1)/2;
                vname[indent] = $2;
                for ( i in vname ) {
                        if (i > indent) {
			delete vname[i]
			}
                }
                if (length($3) > 0) {
                        vn="";
                        for (i=0; i<indent; i++) {
                                vn=(vn)(vname[i])("_")
                        }
                printf("%s%s%s=\"%s\"\n", "'$prefix'",vn, $2, $3);
                }
        }'
}

function read_yml()
{
        eval $(read_config /home/.travis.yml "yml_")
}

function set_local_repo()
{
        color_printlin cyan "Start updating local repository"
        createrepo /home/ > /dev/null
        check_status "Update local repository succeeded" "Update local repository failed"
}

function setup_mock_env()
{
        color_printlin cyan "Start setting up mock environment"
        echo "deltarpm=0" >> /etc/dnf/dnf.conf
        dnf -y install mock rpmdevtools createrepo > /dev/null
        check_status "Setup mock environment succeeded" "Setup mock environment failed"
        color_printlin cyan "Start initializing mock repository"
        /usr/bin/mock -r ${yml_mock_travis_mock_config} --init > /dev/null
        check_status "Mock initialization succeeded" "Mock initialize failed"
}

function main()
{
        read_yml
        setup_mock_env
        add_extra_repo
        build_buildrequires
        set_local_repo
        add_local_repo
        build_target_pkg
}

main
