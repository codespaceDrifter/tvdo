a simple terminal to do app  
clone into a folder. build steps:  
```
go mod init  
go mod tiny
go build tvdo.go
```
the json files will be stored in the same folder  

commands to use:
```
o: create new task
a: edit due date
d: delete task
j: move down task
k: move up task
while editing: enter to send and delete to delete
```
if using linux, put  
alias tvdo="~/yourlocation/tvdo/tvdo"   
in bashrc to call tvdo in terminal  

the root date and default due date is days until 2030. due to AGI maybe arriving 2030. maybe doing tasks past AGI wouldn't make sense for humans. who knows what happens then.  
