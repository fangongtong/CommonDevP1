_G.com = require "tasks.common" 

function RegMe()
	com.Description = Description
	com.TaskConfig = TaskConfig
	com.StartWork = StartWork
	com.SetNewTaskConfig = SetTaskConfig
	--print("regme ok")
end

---------------------------------------------------------------------

description = [[
	这里是文件说明
]]

-- 每个lua文件都应该有CylinderConfig对象用于定义任务对应的缸
--[[CylinderConfig = {
	Cylinder = {
		{Name="left", Pos=nil},
		{Name="right", Pos=nil},
	},
}
]]

LogTaskInfo = {
	Index = 0,
}

local LogTaskPosFields = "PeekForce"

-- 每个lua文件都应该有TaskConfig对象用于定义任务信息
TaskConfig = {
	TaskType = 3,	--1:pull 2:push 3:pull+push
	Force = 0,
	TotalTimes = 0,
	Pos = {0,0}
	PosLimitSw = {0,0,0,0},
}

RunInfo = {
	CurrentTimes = 0,
	ForceDuration = 0,	--ms
	IdleDuration = 0, --ms
	RemainDuration = 0, --s
}

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
	local tstCnt = 100
	local t1 = GetMilliseconds()
	print(t1)
	local t2 = t1
	while  tstCnt > 0 do 
		com.PlcApi_Push_Fix(true, 1,0,0,1000)
		t1 = GetMilliseconds()
		t2 = t1
		while t2-t1 < 1000 do
			com.PlcApi_Nothing(false)
			t2 =  GetMilliseconds()
		end
		
		com.PlcApi_Release(true, 1)
		t1 =  GetMilliseconds()
		t2 = t1
		while t2-t1 < 1000 do
			com.PlcApi_Nothing(false)
			t2 =  GetMilliseconds()
		end
		
		tstCnt = tstCnt -1
	end
	common.PlcApi_Over()
end

function DealPush()
	
end

function dealCheck()
--	print('lua: pos1 force is '..com.RealDt.CR[1].MaxForce)
--	print('lua: pos2 force is '..com.RealDt.CR[2].MaxForce)
--	print('lua: pos3 force is '..com.RealDt.CR[3].MaxForce)
end

function DealPullPush()
	print("DealPullPush in")
	local tstCnt = TaskConfig.TotalTimes
	local t1 = GetMilliseconds()
	local t2 = t1
	while  tstCnt > 0 do 
		com.PlcApi_Push_Fix(true, TaskConfig.Pos[1],0,0,1500)
		com.PlcApi_Push_Fix(true, TaskConfig.Pos[2],0,0,900)
		t1 = GetMilliseconds()
		t2 = t1
		while t2-t1 < 1000 do
			com.PlcApi_Nothing(false)
			t2 =  GetMilliseconds()
		end
		
		dealCheck()
		
		com.PlcApi_Pull_Fix(true, TaskConfig.Pos[1],0,0,1500)
		com.PlcApi_Pull_Fix(true, TaskConfig.Pos[2],0,0,900)
		t1 = GetMilliseconds()
		t2 = t1
		while t2-t1 < 1000 do
			com.PlcApi_Nothing(false)
			t2 =  GetMilliseconds()
		end
		
		dealCheck()
		tstCnt = tstCnt -1
	end
	
	com.PlcApi_Release(true, TaskConfig.Pos[1])
	com.PlcApi_Release(true, TaskConfig.Pos[2])
	print("over")
end

--下面这部分功能可能并不重要，用户可以根据停下来的情况，手动设置次数来继续
function StoreEnv()
	return Common.StoreEnv(TaskConfig,RunInfo)
end

function RestoreEnv(envJson)
	TaskConfig,RunInfo = Common.RestoreEnv(envJson)
end