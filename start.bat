@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

:: UrPicBed Windows 启动脚本

:: 颜色定义
set "RED=[91m"
set "GREEN=[92m"
set "YELLOW=[93m"
set "BLUE=[94m"
set "NC=[0m"

:: 打印带颜色的消息
:print_message
echo %GREEN%[INFO]%NC% %~1
goto :eof

:print_warning
echo %YELLOW%[WARNING]%NC% %~1
goto :eof

:print_error
echo %RED%[ERROR]%NC% %~1
goto :eof

:print_header
echo %BLUE%================================%NC%
echo %BLUE%        UrPicBed 启动脚本        %NC%
echo %BLUE%================================%NC%
goto :eof

:: 检查Docker是否安装
:check_docker
docker --version >nul 2>&1
if errorlevel 1 (
    call :print_error "Docker未安装，请先安装Docker"
    exit /b 1
)

docker-compose --version >nul 2>&1
if errorlevel 1 (
    call :print_error "Docker Compose未安装，请先安装Docker Compose"
    exit /b 1
)
goto :eof

:: 检查配置文件
:check_config
if not exist "config\config.yaml" (
    call :print_error "配置文件 config\config.yaml 不存在"
    exit /b 1
)
call :print_message "配置文件检查通过"
goto :eof

:: 启动服务
:start_service
call :print_message "启动UrPicBed服务..."

:: 停止可能存在的容器
docker-compose down >nul 2>&1

:: 构建并启动服务
docker-compose up -d --build

call :print_message "服务启动完成！"
call :print_message "服务地址: http://localhost:8080"
call :print_message "健康检查: http://localhost:8080/health"
call :print_message "API文档: 请查看README.md"
goto :eof

:: 停止服务
:stop_service
call :print_message "停止UrPicBed服务..."
docker-compose down
call :print_message "服务已停止"
goto :eof

:: 重启服务
:restart_service
call :print_message "重启UrPicBed服务..."
docker-compose restart
call :print_message "服务重启完成"
goto :eof

:: 查看日志
:show_logs
call :print_message "查看服务日志..."
docker-compose logs -f
goto :eof

:: 显示状态
:show_status
call :print_message "服务状态:"
docker-compose ps
goto :eof

:: 显示帮助
:show_help
echo 用法: %~nx0 [命令]
echo.
echo 命令:
echo   start    启动服务
echo   stop     停止服务
echo   restart  重启服务
echo   logs     查看日志
echo   status   查看状态
echo   help     显示此帮助信息
echo.
echo 示例:
echo   %~nx0 start    # 启动服务
echo   %~nx0 logs     # 查看日志
goto :eof

:: 主函数
:main
call :print_header

set "command=%~1"
if "%command%"=="" set "command=start"

if "%command%"=="start" (
    call :check_docker
    if errorlevel 1 exit /b 1
    call :check_config
    if errorlevel 1 exit /b 1
    call :start_service
) else if "%command%"=="stop" (
    call :stop_service
) else if "%command%"=="restart" (
    call :restart_service
) else if "%command%"=="logs" (
    call :show_logs
) else if "%command%"=="status" (
    call :show_status
) else if "%command%"=="help" (
    call :show_help
) else (
    call :print_error "未知命令: %command%"
    call :show_help
    exit /b 1
)

goto :eof

:: 执行主函数
call :main %* 