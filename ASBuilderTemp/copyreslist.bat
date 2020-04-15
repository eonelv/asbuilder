set projectpath=%1
set buildpath=%2
set topath=%3
call "setpath.bat"
copy /y %ASBUILDER_HOME%\build\%projectpath%\%buildpath%\reslist_srv.json %ASBUILDER_HOME%\build\%projectpath%\%topath%