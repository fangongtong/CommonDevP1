
description = [[
	这里是文件说明
	本文件用于单缸的单拉/单压试验
]]


-- 每个lua文件都应该有CylinderConfig对象用于定义任务对应的缸
CylinderConfig = {
	Cylinder = {
		{Name="cylinder", Pos=nil},
	},
}

-- 每个lua文件都应该有TaskConfig对象用于定义任务信息
TaskConfig = {
	TaskType = 0,	--1:pull 2:push
	Force = 0,
	Frequency = 0,
	TotalTimes = 0,
}

RunInfo = {
	CurrentTimes = 0,
	ForceDuration = 0,	--ms
	IdleDuration = 0, --ms
	RemainDuration = 0, --s
}

Common = require "common.lua" 

--function initCylinder(cyl)
--	print(cyl.Name)
--end

function StopWork()
	print("StopWork")
	
	
	return 0
end


function CheckConfig()
	if Common.CheckCylinderConfig(CylinderConfig) == false then
		return Common.Err_CylinderConfig
	end
	
	if not (TaskConfig.TaskType > 0 and TaskConfig.TaskType < 3) then 
		return Common.Err_TaskConfig
	end
	
	if Force<1 or Force > 10000 then 
		return Common.Err_TaskConfig
	end
	
	if Frequency < 1 or Frequency > 5 then 
		return Common.Err_TaskConfig
	end
	
	if TotalTimes < 1 or TotalTimes > 20000 then 
		return Common.Err_TaskConfig
	end
	
	return 0
end

function InitCylinder()
	RunInfo.IdleDuration = ((1000 / TaskConfig.Frequency) - Common.CommuTime*2)/2
	RunInfo.ForceDuration = RunInfo.IdleDuration+Common.CommuTime
end

function StartWork()
	print("StartWork")
	
	while  RunInfo.CurrentTimes < TaskConfig.TotalTimes do 
		-- call go code 发力
		GoSleep(RunInfo.ForceDuration)
		-- call go code 卸力
		GoSleep(RunInfo.IdleDuration)
	end
	
	return 0
end


--下面这部分功能可能并不重要，用户可以根据停下来的情况，手动设置次数来继续
function StoreEnv()
	return Common.StoreEnv(TaskConfig,RunInfo)
end

function RestoreEnv(envJson)
	TaskConfig,RunInfo = Common.RestoreEnv(envJson)
end