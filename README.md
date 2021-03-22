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
FmsMasterURL|fmshost1:4544,fmshost2:4544|Comma separated list with the master FMSs where the app will look.

## Usage

- TODO

## Endpoints

### serverhealth
    200 OK
    http(s)://servername:port/serverhealt
    
    Shows if this server is running

### fscstatus
    200 OK
    http(s)://servername:port/fscstatus/

    Returns a JSON array of FSC Status objects with the status of all FSC configured in the FMSMaster(s) configured in settings.ini

### fscstatus/$host
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

### fscalive/$host
    200 OK
    http(s)://servername:port/fscalive/$fschostname
    
    Returns a JSON if the fsc server is alive, for example: http://fmshelper.yolo.com:8080/fscalive/foo-bar-fsc1 (optional query parameter port: ?port=1234)

```
200 OK
{"status":"OK"}
```

### fscversion/$host
    200 OK
    http(s)://servername:port/fscversion/$fschostname
    
    Returns a JSON with the components version of the passed FSC, for example: http://fmshelper.yolo.com:8080/fscversion/foo-bar-fsc1 (optional query parameter port: ?port=1234)

```
200 OK
{
   "fms_server_cache":{
      "version":"11.6.0",
      "build_date":"20190930"
   },
   "fms_util":{
      "version":"11.6.0",
      "build_date":"20190930"
   },
   "fsc_java_client_proxy":{
      "version":"11.6.0",
      "build_date":"20190930"
   }
}
```
### fscconfig
    200 OK
    http(s)://servername:port/fscconfig/

    Returns the native TC XML document with the FMS configuration from the declared FMS master servers in settings.ini. 
    It will look for the configuration in order and returns the first server that responds to the command without errors.
### fscconfig/$host
    200 OK
    http(s)://servername:port/fscconfig/$host

    Returns the native TC XML document with the FMS configuration from the passed FMS/FSC server, for example: http://fmshelper.yolo.com:8080/fscconfig/foo-bar-fsc1 (optional query parameter port: ?port=1234)

### fscconfig/$host/hash/
    200 OK
    http(s)://servername:port/fscconfig/$host/hash

    Returns a string with the FMS configuration hash from the passed FMS/FSC server, for example: http://fmshelper.yolo.com:8080/fscconfig/foo-bar-fsc1/hash/ (optional query parameter port: ?port=1234)

### fmsconfigreport/
    200 OK
    http(s)://servername:port/fscconfigreport/

    Returns a JSON with an array of objects with the result of the FSCAdmin ./config/report command for the first available configured FMSMaster in settings.ini

```json
[
    {
      "status":"OK",
      "fsc_id":"FMS01",
      "is_master":true, //True for FMS master servers, false for slaves
      "config_hash":"558c95a37eaee435e45a0e4e1400f097"
   },
   {
      "status":"OK",
      "fsc_id":"FSC01",
      "is_master":false,
      "config_hash":"558c95a37eaee435e45a0e4e1400f097" //Only shown if not error is thrown
   },
   {
      "status":"KO",
      "fsc_id":"FSC02",
      "is_master":false,
      "error":"ERROR_ALL_LINKS_DOWN_1{[[-1946043301^FMS01 -1945380301^AGROUP -194653401^FSC02]]}" //The error is only shown when exists
   }
]
```

### log/$host

    200 OK
    http(s)://servername:port/log/$host

    Returns the plain current text log from the passed FMS/FSC server, for example: http://fmshelper.yolo.com:8080/log/foo-bar-fsc1 
    - Optional query parameter "port": ?port=1234
    - Optional query parameter "lines": ?lines=100 (valid values: "all" or a valid integer, for example: 100)


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
