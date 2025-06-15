_G.com = require "tasks.common" 

function RegMe()
	com.Description = Description
	com.TaskConfig = TaskConfig
	com.StartWork = StartWork
	com.SetNewTaskConfig = SetTaskConfig
	print("regme ok")
end

Description = [[
	这里是文件说明
	本文件用于单缸的单拉/单压试验
]]

-- 每个lua文件都应该有TaskConfig对象用于定义任务信息
TaskConfig = {
	TaskType = 3,	--1:pull 2:push 3:pull+push
	Force = 0,
	Frequency = 0,
	TotalTimes = 0,
	Pos = {0,0},
	PosLimitSw = {0,0,0,0},
}


function SetTaskConfig(tskCfg)
	TaskConfig.TaskType = tskCfg.TaskType
	TaskConfig.Force = tskCfg.Force
	TaskConfig.TotalTimes = tskCfg.TotalTimes
	--print(tskCfg.abc.a1.b[1].b2)
end
--[[
RunInfo = {
	CurrentTimes = 0,
	ForceDuration = 0,	--ms
	IdleDuration = 0, --ms
	RemainDuration = 0, --s
}
]]


--local socket=require'socket'

--function initCylinder(cyl)
--	print(cyl.Name)
--end

function StopWork()
	print("StopWork")
	
	return 0
end


function CheckConfig()
--	if Common.CheckCylinderConfig(CylinderConfig) == false then
--		return Common.Err_CylinderConfig
--	end
	
--	if not (TaskConfig.TaskType > 0 and TaskConfig.TaskType < 3) then 
--		return Common.Err_TaskConfig
--	end
	
--	if Force<1 or Force > 10000 then 
--		return Common.Err_TaskConfig
--	end
	
--	if Frequency < 1 or Frequency > 5 then 
--		return Common.Err_TaskConfig
--	end
	
--	if TotalTimes < 1 or TotalTimes > 20000 then 
--		return Common.Err_TaskConfig
--	end
	
	return 0
end

function InitCylinder()
--	RunInfo.IdleDuration = (1000 / TaskConfig.Frequency) - Common.CommuTime*2)/2
--	RunInfo.ForceDuration = RunInfo.IdleDuration+Common.CommuTime
end

function StartWork()
	print("StartWork")
	
	if TaskConfig.TaskType == 1 then
	print("StartWork 1")
		DealPull()
	elseif TaskConfig.TaskType == 2 then
	print("StartWork 2")
		
	elseif TaskConfig.TaskType == 3 then
	print("StartWork 3")
		DealPullPush()
	else
	print("StartWork ..")
	end
	
	
	
	--[[
	while  RunInfo.CurrentTimes < TaskConfig.TotalTimes do 
		-- call go code 发力
		GoSleep(RunInfo.ForceDuration)
		-- call go code 卸力
		GoSleep(RunInfo.IdleDuration)
	end
	]]
	
	com.PlcApi_Over()
end

function DealPull()
	print("DealPull in")
	--local tstCnt = 100
	print("DealPull GetMilliseconds")
	local t1 = GetMilliseconds()
	print(t1)
	local t2 = t1
	while  TaskConfig.TotalTimes > 0 do 
	
		com.PlcApi_Push_Fix(true, TaskConfig.Pos[1],0,0,TaskConfig.Force)
		t1 = GetMilliseconds()
		t2 = t1
		while t2-t1 < 1000 do
		
			com.PlcApi_Nothing(false)
			t2 =  GetMilliseconds()
		end
		
	print("DealPull 3")
		com.PlcApi_Release(true, TaskConfig.Pos[1])
		t1 =  GetMilliseconds()
		t2 = t1
		while t2-t1 < 1000 do
	
			com.PlcApi_Nothing(false)
			t2 =  GetMilliseconds()
		end
		
		TaskConfig.TotalTimes = TaskConfig.TotalTimes -1
	end
end

function DealPush()
	
end

function DealPullPush()
	print("DealPullPush in")
	local t1 = GetMilliseconds()
	print(t1)
	local t2 = t1
	while  TaskConfig.TotalTimes > 0 do 
		com.PlcApi_Push_Fix(true, TaskConfig.Pos[1],0,0,TaskConfig.Force)
		t1 = GetMilliseconds()
		t2 = t1
		while t2-t1 < 1000 do
			com.PlcApi_Nothing(false)
			t2 =  GetMilliseconds()
		end
		
		
		com.PlcApi_Pull_Fix(true, TaskConfig.Pos[1],0,0,TaskConfig.Force)
		t1 = GetMilliseconds()
		t2 = t1
		while t2-t1 < 1000 do
			com.PlcApi_Nothing(false)
			t2 =  GetMilliseconds()
		end
		
		TaskConfig.TotalTimes = TaskConfig.TotalTimes -1
	end
	
	com.PlcApi_Release(true, TaskConfig.Pos[1])
	print("over")
end

--下面这部分功能可能并不重要，用户可以根据停下来的情况，手动设置次数来继续
function StoreEnv()
	return Common.StoreEnv(TaskConfig,RunInfo)
end

function RestoreEnv(envJson)
	TaskConfig,RunInfo = Common.RestoreEnv(envJson)
end