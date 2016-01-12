# GoIsilon

## Overview
```GoIsilon``` represents API bindings for Go that allow you to manage Isilon
NAS platforms.  In the true nature of API bindings, it is intended that the
functions available are basically a direct implementation of what is available
through the API.  There is however, an abstraction called object which can be
used by things like `Docker`, `Mesos`, and
[REX-Ray](https://github.com/emccode/rexray) to integrate Isilon.

## API Compatibility
Currently only tested with v7+.

## Examples
The package was written using test files, so these can be looked at for a more
comprehensive view of how to implement the different functions.

Intialize a new client

	c, err := NewClient() // or NewClientWithArgs(endpoint, insecure, userName,
    password,volumePath)
	if err != nil {
		panic(err)
	}


Create a Volume

    volume, err := c.CreateVolume("testing")


Export a Volume

    err := c.ExportVolume("testing")


Delete a Volume

    _, err := c.DeleteVolume(name)



For example usage you can see the [REX-Ray](https://github.com/emccode/rexray)
repo.  There, the ```goisilon``` package is used to implement a
```Volume Manager``` across multiple storage platforms. This includes managing
multipathing, mounts, and filesystems.

## Environment Variables
Name | Description
---- | -----------
`GOISILON_ENDPOINT` | the API endpoint, https://172.17.177.230:8080
`GOISILON_USERNAME` | the username
`GOISILON_GROUP` | the user's group
`GOISILON_PASSWORD` | the password
`GOISILON_INSECURE` | whether to skip SSL validation
`GOISILON_VOLUMEPATH` | which base path to use when looking for volume directories

## Contributions
Please contribute!

Licensing
---------
Licensed under the Apache License, Version 2.0 (the “License”); you may not use
this file except in compliance with the License. You may obtain a copy of the
License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an “AS IS” BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.

Support
-------
If you have questions relating to the project, please either post
[Github Issues](https://github.com/emccode/mesos-module-dvdi/issues), join our
Slack channel available by signup through
[community.emc.com](https://community.emccode.com) and post questions into
`#projects`, or reach out to the maintainers directly.  The code and
documentation are released with no warranties or SLAs and are intended to be
supported through a community driven process.
