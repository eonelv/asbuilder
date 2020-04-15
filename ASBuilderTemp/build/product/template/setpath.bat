REM command is project ID
REM svn://192.168.0.16/client/appoint-gods/trunk/appoint-gods
set SVN_MAIN=@code/appoint-gods/src
set SVN_ELIB=@code/e3de/source/src
set SVN_RES=@res/out
set SVN_UIEDITOR= @code/ag
set SVN_swcs= svn://192.168.0.16/client/libs

call "..\..\..\setpath.bat"
set DIRPATH=%ASBUILDER_HOME%\build\%projectpath%\%basepath%\