package consts

// Key to use when setting the user info
type ctxKeySystemUserInfo int

const (
	LOGO string = ` __   __               _____                     
 \ \ / /              |  __ \                    
  \ V / __ _ _ __ _ __| |__) | __ _____  ___   _ 
   > < / _` + "`" + ` | '__| '__|  ___/ '__/ _ \ \/ / | | |
  / . \ (_| | |  | |  | |   | | | (_) >  <| |_| |
 /_/ \_\__,_|_|  |_|  |_|   |_|  \___/_/\_\\__, |
                                            __/ |
                                           |___/ 

%s Copyright (c) 2023-2023 Build with ❤️ By d0zingcat
`

	VERSION = "v0.1.0"
	AUTHOR  = "d0zingcat"
	REPO    = "xarr-proxy"

	SQL_FILE_DIR        = "resources/sql"
	CHECKPOINT_FILENAME = "CHECKPOINT"
	STATIC_FILE_DIR     = "resources/static"
	RULE_FILE_DIR       = "resources/rule"

	USER_INFO_CTX_KEY ctxKeySystemUserInfo = 0
)
