#!/bin/bash

# Global variables
SCRIPT_URL="https://raw.githubusercontent.com/nrechn/mock-travis/master/"
SCRIPT="mock.sh"
LOCATION=`pwd`

# Variables for docker
DOCKER_OPTION="--cap-add=SYS_ADMIN --privileged=true"


function check_overall_status()
{
        if [ $? != 0 ]
        then {
                echo -e "\n\e[31m\e[1mOVERALL: Fail to build $1 and related build dependencies.\e[0m"
        } else {
                echo -e "\n\e[33m\e[1mOVERALL: Successfully build $1 and related build dependencies.\e[0m"
        }
        fi
}

function prep_scrip()
{
        cd ${LOCATION}/
        wget ${SCRIPT_URL}${SCRIPT} > /dev/null 2>&1
        chmod +x ${SCRIPT}
}

function pull_docker()
{
        docker pull $1
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
        eval $(read_config .travis.yml "yml_")
}

function remove_docker_container()
{
        docker rm mock-build > /dev/null 2>&1
}

function run_docker()
{
        IMAGE=${yml_mock_travis_docker_image}
        docker run --name mock-build ${DOCKER_OPTION} -v ${LOCATION}/:/home -i ${IMAGE} /home/${SCRIPT}
}

function main()
{
        read_yml
        prep_scrip
        pull_docker ${yml_mock_travis_docker_image}
        run_docker
        check_overall_status ${yml_mock_travis_packages_name}
        remove_docker_container
}

main
