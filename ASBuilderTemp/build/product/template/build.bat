set command=%1
set projectpath=%2
set basepath=%3
set patch=%4
set isupdate=%5
cd build\%projectpath%\%basepath%
call "buildmain.bat"