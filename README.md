# Mock-Travis
Mock-Travis utilizes [Travis CI](https://travis-ci.org/) to provide a continuous integration practice for testing **`spec`** files of RedHat RPM packages. That is to say, Mock-Travis can be utilized to build source rpm packages and binary rpm packages in order to check whether **`spec`** files are written correctly.

> **Note**: RPM package **spec** files need to be stored in the GitHub repository.

<br>
### How to use Mock-Travis?
Setting up Mock-Travis is quite simple. All you need to do are just two things:
- Put **`.travis.yml`** your GitHub repository.
- Enable autobuild on [Travis CI](https://travis-ci.org/) website.

> Assumption: You should know how to add file and push to your GitHub repository; how to [sign in to Travis CI](https://travis-ci.org/auth) with your GitHub account, and go to [profile page](https://travis-ci.org/profile) and enable [Travis CI](https://travis-ci.org/) for the repository you want to build.

<br>
**Example of `.travis.yml` file**
```yml
#! -------------------------------------------------------------------------
#! Mock-Travis Settings
#! -------------------------------------------------------------------------

mock_travis:
  # Set docker image which will be utilized.
  docker_image: fedora

  # Set a particular test configuration for mock.
  mock_config: fedora-23-x86_64

  # Set the packages to be tested.
  # Packages directory.
  # For example
  #         If the PACKAGE spec file is stored like this
  #         ├───GitHub Repository
  #         │   └───rpm
  #         │       └───build
  #         │           └───PACKAGE
  #         │               └───PACKAGE.spec
  #         The directory should be set to `./rpm/build` or `rpm/build`.
  #
  #         If the PACKAGE spec file is stored like this
  #         ├───GitHub Repository
  #         │   └───PACKAGE
  #         │       └───PACKAGE.spec
  #         The directory should be set to `./` or leave it empty.
  packages_directory: ./
  # Target packages.
  # It is allowed to set multiple package name in one line.
  packages_name: 
  # Build requires of target package.
  # It is allowed to set multiple package name in one line.
  # Set here will force script to ignore the `git` below.
  packages_buildrequires: 
  # Use RPM spec files from a GitHub repository
  # instead of storing in local.
  # It only works if the `buildrequires` above is empty.
  # Otherwise, script will use `buildrequires` and ignore the `git`.
  # For example
  #         If your spec repository is 
  #         "https://github.com/nrechn/Sway-Fedora",
  #         git could be set to `nrechn/Sway-Fedora`
  packages_buildrequires_git: 



#! -------------------------------------------------------------------------
#! DO NOT EDIT THE FOLLOWING SETTINGS
#! UNLESS YOU KNOW WHAT YOU ARE DOING
#! -------------------------------------------------------------------------

sudo: required

services:
  - docker

script:
  - bash -c "$(wget -q https://raw.githubusercontent.com/nrechn/mock-travis/master/build.sh -O -)"
```

You can simply copy and paste the example above, or download the [example.travis.yml](https://raw.githubusercontent.com/nrechn/mock-travis/master/example.travis.yml).

> **Note**: Please remember to rename **[example.travis.yml](https://raw.githubusercontent.com/nrechn/mock-travis/master/example.travis.yml)** to **.travis.yml** if you download **[example.travis.yml](https://raw.githubusercontent.com/nrechn/mock-travis/master/example.travis.yml)**.

<br>
### How Mock-Travis works?
When you make a push to your GitHub repository, it will trigger a [Travis CI](https://travis-ci.org/) build. The build process will run a docker container and do the following things:
- Install necessary tools. (e.g. `mock`, `spectool`, etc.)
- Initialize mock config.
- Build dependencies source packages.
- Build dependencies binary packages.
- Create local repository and add it to the mock config.
- Build target source packages.
- Build target binary packages.

The [Travis CI](https://travis-ci.org/) will show the whole build log. The results of each step can be found easily in build log as they will be shown in colored bold words.

[nrechn/Sway-Fedora](https://github.com/nrechn/Sway-Fedora) is a real world example of utilizing Mock-Travis. You can check its [.travis.yml](https://github.com/nrechn/Sway-Fedora/blob/master/.travis.yml) file, or [Travis CI build log](https://travis-ci.org/nrechn/Sway-Fedora).

<br>
### Limitations
- The customizability is still low. Only provide a few options currently.
- Run `mock` in docker container to test packages is few minutes slower than test locally.
