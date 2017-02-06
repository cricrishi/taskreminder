Software Requirments 

Golang (1.5)
Mysql

----------------------------------------------------------------

Create the mysql table

CREATE TABLE IF NOT EXISTS `task`(
  `task_id` int(11) NOT NULL AUTO_INCREMENT,
  `mobile_number` VARCHAR(12) NOT NULL,
  `email` VARCHAR(50) NOT NULL,
  `reminder_message` VARCHAR(500) NOT NULL,
  `reminder_time`  datetime NOT NULL,
  `reminder_status` boolean NOT NULL DEFAULT 0,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created Time',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated Time',
   PRIMARY KEY (`task_id`)
 ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


To run type go run main.go

Apis End Url :- 

Get /tasks

Post /task/create

Jsoon Data :- 

{
"email" : "rishi",
"mobile"" "8800381831"
"reminder_message" : "test"
"reminder_time": "2016-08-15 16:56:00"
"reminder_text" : "this is first text" 
}
hello this is rebasing
