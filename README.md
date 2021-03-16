# tcfmshelper

tcfmshelper is a REST service built ontop of Teamcenter FMS/FSC **fscadmin** command to help Teamcenter administrators to query FSC server status. 

It works by invoking using  ```os\exec``` the default FSCAdmin java command. The configuration loaded at the begining of the server lauch process mimics the fscadmin.bat/sh variable and invokation steps.

## Installation

TODO

## Configuration

All the configuration is done through the conf/settings.ini file and standard TC environment variables

### Environment Variables:

- ```JAVA_HOME``` should point to a valid Java JRE folder. The proccess will look for ```bin/java``` inside this folder. This is **mandatory**.

- ```FMS_HOME``` should point to a valid folder containing fscadmin.bat/sh script and **jar** folder with the Teamcenter FSC jars, for example: ```FMS_HOME=C:\Siemens\TC116\fsc```. This is **optional** as it could be defined in ```settings.ini``` also.

### settings.ini

**Name**|**Value**|**Comments**
-----|-----|-----
**[server]**| | 
RunMode|release|Valid values: debug, release
Port|8080| 
ReadTimeout|60s|With unit (s for seconds, m for minutes)
WriteTimeout|60s|With unit (s for seconds, m for minutes)
**[fsc]**| | 
FscLocation|C:\Siemens\TC116\fsc|Not mandatory. It can be specified with the Environment Variable FMS\_HOME. Overrides the variable if specified.
## Usage

- TODO

## Endpoints

### serverhealth
    200 OK
    http(s)://servername:port/serverhealt
    
    Shows if this server is running

### fscstatus
    200 OK
    http(s)://servername:port/fscstatus/$fschostname
    
    Returns a JSON with the status of the passed FSC, for example: http://fmshelper.yolo.com:8080/fscstatus/foo-bar-fsc1 (optional query parameter port: ?port=1234)

```
200 OK
{
    "status":"OK",
    "fsc_id":"FSC_FOO01",
    "site":"-194608575",
    "current_admin_connections":1,
    "current_file_connections": 6
}
```

### fscalive
    200 OK
    http(s)://servername:port/fscalive/$fschostname
    
    Returns a JSON if the fsc server is alive, for example: http://fmshelper.yolo.com:8080/fscalive/foo-bar-fsc1 (optional query parameter port: ?port=1234)

```
200 OK
{"status":"OK"}
```

### Error handling 

Right now, in this development stage, all errors are returned with a 500 server error and a JSON object like so:

```
500 SERVER ERROR
{
    "status":"KO",
    "message":"Error message: reason"
}
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)
