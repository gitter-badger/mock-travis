# Mock-Travis

[![Build Status](http://img.shields.io/travis/nrechn/mock-travis.svg?style=flat-square)](https://travis-ci.org/nrechn/mock-travis)

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
  # Set a particular test configuration for mock.
  mock_config: fedora-23-x86_64

  # Use RPM spec files from a GitHub repository
  # instead of storing in local.
  # It only works if the `buildrequires` above is empty.
  # Otherwise, script will use `buildrequires` and ignore the `git`.
  # For example
  #         If your spec repository is 
  #         "https://github.com/nrechn/Sway-Fedora",
  #         git could be set to `nrechn/Sway-Fedora`
  packages_buildrequires_git: 
  # Use extra/external repository during building packages.
  # This option allows mock to access an additional repository
  # plus defaults repositories based on the mock config.
  # gpgcheck is disabled in this option.
  # For example
  #         If add FZUG as an extra repository
  #         packages_extra_repo should set to
  #         "https://repo.fdzh.org/FZUG/testing/24/x86_64/"
  packages_extra_repo:


#! -------------------------------------------------------------------------
#! DO NOT EDIT THE FOLLOWING SETTINGS
#! UNLESS YOU KNOW WHAT YOU ARE DOING
#! -------------------------------------------------------------------------

sudo: required

services:
  - docker

script:
  - wget -q https://github.com/nrechn/mock-travis/releases/download/latest/mock-travis 
        && chmod +x mock-travis 
        && ./mock-travis
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

<br>
### Advantages
- No need to test mock build on your own computer. It is quite hard to run mock if you use other GNU/Linux distros than RedHat related GNU/Linux distros.
- Colorful output. The [Travis CI](https://travis-ci.org/) will show the whole build log. The results of each step can be found easily in build log as they will be shown in colored bold words.

Here is an example of build log:
![Travis-CI log](https://github.com/nrechn/mock-travis/raw/master/misc/travis-ci-log.png)

> Click the picture to view the raw file.

<br>
### Projects use Mock-Travis
[nrechn/Sway-Fedora](https://github.com/nrechn/Sway-Fedora) is a real world example of utilizing Mock-Travis. You can check its [.travis.yml](https://github.com/nrechn/Sway-Fedora/blob/master/.travis.yml) file, or [Travis CI build log](https://travis-ci.org/nrechn/Sway-Fedora).


<br>
### Limitations
- The customizability is still low. Only provide a few options currently.
- Run `mock` in docker container to test packages is few minutes slower than test locally.
