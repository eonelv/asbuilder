REM command is project ID
set command=1
set SVN_LIB=svn://192.168.0.10/eonegame/project_flash/trunk/e3d-engine/library
set SVN_MAIN=svn://192.168.0.10/eonegame/project_flash/trunk/e3d-engine/betray-gods/src
set SVN_ELIB=svn://192.168.0.10/eonegame/project_flash/trunk/e3d-engine/e3de/source/src
set SVN_RES=svn://192.168.0.10/eonegame/project_flash/trunk/e3d-engine/resources/out
set SVN_UIEDITOR= svn://192.168.0.10/eonegame/project_flash/trunk/e3d-engine/resources/ui/project/exports

call "..\..\..\setpath.bat"
set DIRPATH=%ASBUILDER_HOME%\build\%projectpath%\%basepath%\