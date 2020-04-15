@echo ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
@echo ~  eone - AS3 project builder
@echo ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

@echo off
REM 初始化所有环境变量
call "setpath.bat"

cd \
cd %DIRPATH%

echo SysBuildMsg=Begin check out code and resource
:ck
if not exist "source" (md "source")
if not exist "source"/src (svn checkout %SVN_MAIN% "source"/src)
if not exist "source"/lib (svn checkout %SVN_LIB% "source"/lib)
if not exist "source"/elib (svn checkout %SVN_ELIB% "source"/elib)
if not exist "source"/res (svn checkout %SVN_RES% "source"/res)
if not exist "source"/ag (svn checkout %SVN_UIEDITOR% "source"/ag)
if not exist "swcs" (svn checkout %SVN_swcs% "swcs")

:up
if %isupdate% == 1 (
@echo %basepath% 目录存在, 准备更新代码
svn cleanup source/src
svn cleanup source/ag/e1/e2d/gen
svn cleanup source/elib
svn cleanup source/res
svn update source/src > log_code.txt --accept tc
svn update source/lib >> log_code.txt --accept tc
svn update source/elib >> log_code.txt --accept tc
svn update source/ag/e1/e2d/gen >> log_code.txt --accept tc
svn update source/res > log_res.txt --accept tc
svn update swcs > log_lib.txt --accept tc
)

if exist "release" (RD "release" /S /Q)
java -Xms512m -Xmx512m e1.tools.asbuilder.ASBuiler %DIRPATH%\"source" %ASBUILDER_HOME% %command% %patch% %isupdate% %DIRPATH% %DIRPATH%source\res %DIRPATH%deploy\release %DIRPATH%

REM "C:\Program Files (x86)\JsonAMF\JsonAMF.exe" %DIRPATH%deploy\release
REM java e1.tools.asbuilder.CreateAMF %DIRPATH%deploy\release

native2ascii -encoding utf-8 build.properties temp.properties

if not exist "log" (md "log")
:compiler
@echo.
@echo 开始编译，请耐心等待...
@echo 编译后产生log文件，请查看log\eone-betraygods-%basepath%.log
ant -buildfile build.xml -logfile log\eone-betraygods-%basepath%.log

:end
echo 用户退出