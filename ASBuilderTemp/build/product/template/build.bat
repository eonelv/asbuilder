set projectpath=%1
set basepath=%2
set patch=%3
set isupdate=%4
cd build\%projectpath%\%basepath%
call "buildmain.bat"