# NAS System

## Introduction
This is a NAS (Network Attached Storage) system that can be deployed on PC after compile.

#### _Disclaimers_
This is just a personal project, so we are not responsible for your data security.
We cannot guarantee that you will not encounter bugs causing data errors.
However, we welcome for any report on bugs. Thank you for your contributions.

## Installation

1. use `git clone` or simply download and decompression to get the whole project.  
**_Back-end_**
2. prepare certificates for https connection.
3. modify the config.yml to specify features of back-end according to those comments.
4. double click to execute nas_server. (you can also execute it in terminal)
5. visit `https://ip_of_host:server_port/hello` to test the back-end. The default port is 443.
6. compile....  
**_Front-end_**
7. xxx

## Instruction

1.  visit `https://ip_of_host:front_end_port/xxx` to start your journey.
2. 

## Architecture
This project is of typical MVC style architecture. 

### Front-end
Vue is used as the framework of the front-end.

### Back-end
The whole back-end can be roughly divided to 3 layers. 
They are controller layer, service layer and data access layer.
- Controller Layer  
This layer is responsible for receiving request.
- Service Layer  
This layer is responsible for handling some specific tasks.
- Data Access Layer  
This layer is responsible for the file access.

## The principle of file transfer
### Small files
### Large files
- Upload:  
- Download

## Contributors
- Minghao Han --Responsible for architecture design and back-end development.
- Jinhui Huang --Responsible for front-end development.