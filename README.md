
# Oimo Admin
This plugin can help you quickly create a admin with only one line code to implement CRUD functions, 
which is very suitable for small and simple projects. 

## Install
`go get -u github.com/oimoyu/OimoAdmin@latest`

## Code Usage
*Set up your gin router and gorm db conn, then pass them to `OimoAdmin.Init` function. After running gin, you will see the login info in the console output.*
```
    var DB *gorm.DB
    var r *gin.Engine
    
    // init gin router and gorm conn
    // ...
    
    OimoAdmin.Init(r, DB)
    r.Run("0.0.0.0:8098")
```
You can check the full example code [here](./example/example.go). 
The full example code needs sqlite driver `go get -u gorm.io/glebarez/sqlite`, you can use other database driver if you want.


## Usage Notice

### ID Field Needed
Update and delete operations rely on the `id` field. When you create table using gorm, please set `id` field as primary key

gorm table reference settings:
```
type Order struct {
    ID         int `gorm:"primaryKey" json:"id"`
    ...
}
```

## Runtime File
After you run the code, this plugin will create a directory named `oimo_admin_runtime` in the program directory to store runtime files.

### ①Frontend Requirement
This admin rely on [amis](https://github.com/baidu/amis) to build the admin UI. Plugin will automatically download the sdk.
If you have network issue, you can [download it](https://github.com/baidu/amis/releases/download/6.6.0/sdk.tar.gz) 
manually and unzip the file to `oimo_admin_runtime/wwwroot/sdk` so that `sdk.js` is directly under the `sdk` folder.

### ②Environment Files
`oimo_admin_runtime/.env` stores the admin key file and log file. 
If you want to reset the admin login key, you can delete the secret files and re-run the code

## Reverse Proxy
This plugin needs the client's real IP to ensure that fail2ban and logs are recorded normally. 
If you want to run behind a reverse proxy, please set the header correctly.

nginx reference settings:
```
    location / {
        proxy_pass http://127.0.0.1:8098;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forward-For $remote_addr;
    }
```

## Safety
The admin login of this plugin uses multiple checks to ensure security.

Please do not expose the .env folder to outside, as it contains secret files.

① random username and password: avoid weak username and password.

② random path: using random path to hide the api access.

③ login IP check: check last login ip.

④ single active token: check last login token.

⑤ fail to ban ip: when a IP tries to log in multiple times with incorrect credentials, the IP will be banned,
restart to clear banned ip.

