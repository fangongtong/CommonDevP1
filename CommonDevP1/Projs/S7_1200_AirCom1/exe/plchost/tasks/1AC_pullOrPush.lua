_G.com = require "tasks.common" 

function RegMe()
	com.Description = Description
	com.TaskConfig = TaskConfig
	com.StartWork = StartWork
	com.SetNewTaskConfig = SetTaskConfig
	--print("regme ok")
end

---------------------------------------------------------------------

Description = [[
	这里是文件说明
	本文件用于单缸的单拉/单压试验
]]

LogTaskInfo = {
	Index = 0,
}

local LogTaskPosFields = "PeekForce"

-- 每个lua文件都应该有TaskConfig对象用于定义任务信息
TaskConfig = {
	TaskType = 3,	--1:pull 2:push 3:pull+push
	Force = 0,
	Frequency = 0,
	TotalTimes = 0,
	Pos = {0},
	PosLimitSw = {0,0},
}


function SetTaskConfig(tskCfg)
	TaskConfig.TaskType = tskCfg.TaskType
	TaskConfig.Force = tskCfg.Force
	TaskConfig.TotalTimes = tskCfg.TotalTimes
	--print(tskCfg.abc.a1.b[1].b2)
end


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
		DealPull()
	elseif TaskConfig.TaskType == 2 then
		DealPush()
	elseif TaskConfig.TaskType == 3 then
		DealPullPush()
	else
		Test5()
	end
	
	com.PlcApi_Over()
end

function Test2()
	if com.LoadMenu(TaskConfig.Pos[1]) == true then
		print(com.MenuConfig[TaskConfig.Pos[1]].OrderMaxRunSeconds           )
		print(com.MenuConfig[TaskConfig.Pos[1]].ForceSensorPushOverloadRate  )
		print(com.MenuConfig[TaskConfig.Pos[1]].ForceSensorPullOverloadRate  )
		print(com.MenuConfig[TaskConfig.Pos[1]].SampleForceOverloadRate      )
		print(com.MenuConfig[TaskConfig.Pos[1]].ForceFallProtectRate         )
		print(com.MenuConfig[TaskConfig.Pos[1]].MaxSecondsBeforeReachForce   )
		print(com.MenuConfig[TaskConfig.Pos[1]].ProportionalValve            )
		print(com.MenuConfig[TaskConfig.Pos[1]].TargetForceUpperDeviationRate)
		print(com.MenuConfig[TaskConfig.Pos[1]].TargetForceLowerDeviationRate)
		print(com.MenuConfig[TaskConfig.Pos[1]].ForceFactorB                 )
		print(com.MenuConfig[TaskConfig.Pos[1]].ForceSensorCapacity          )
	end
	
	local t1 = GetMilliseconds()
	for i = 1,100 do
		com.PlcApi_Set_GainRate(false,5, i/10 )
	end
	t1 = GetMilliseconds()-t1
	print("used time:"..t1)
end


function Test1()
	local t1 = GetMilliseconds()
	local t2 = t1
	local adjustRange = 100
	local adjust = 800
	while  adjust < 3000 do 
		adjust = adjust + adjustRange
		com.PlcApi_Push_Fix(true, TaskConfig.Pos[1],0,0, adjust)
		t1 = GetMilliseconds()
		t2 = t1
		while t2-t1 < 4000 do
			com.PlcApi_Nothing(false)
			t2 =  GetMilliseconds()
		end
		print(string.format("adj: %d RForce: %.2f MForce: %.2f", adjust, com.RealDt.CR[ TaskConfig.Pos[1] ].RealForce,com.RealDt.CR[ TaskConfig.Pos[1] ].MaxForce))
	end
	
	com.PlcApi_Release(true, TaskConfig.Pos[1])
end

function Test5()
	local t1 = GetMilliseconds()
	local t2 = t1
	local cnt = 0
	
	while  cnt < 20 do 
		t1 = GetMilliseconds()
		t2 = GetMilliseconds()
--		if t2-t1 > 2 then
--			print("too long: %d", t2-t1)
--		end
			print("too long: %d", t2-t1)
		com.PlcApi_Nothing(false)
		cnt = cnt+1
	end
	
end

function DealPull()
	local tstCnt = 100
	local t1 = GetMilliseconds()
	local t2 = t1
	
	while  tstCnt > 0 do 
		com.PlcApi_Push_Fix(true, TaskConfig.Pos[1],0,0,1000)
		t1 = GetMilliseconds()
		t2 = t1
		while t2-t1 < 1000 do
			com.PlcApi_Nothing(false)
			t2 =  GetMilliseconds()
		end
		
		com.PlcApi_Release(true, TaskConfig.Pos[1])
		t1 =  GetMilliseconds()
		t2 = t1
		while t2-t1 < 1000 do
			com.PlcApi_Nothing(false)
			t2 =  GetMilliseconds()
		end
		
		tstCnt = tstCnt -1
	end
end

function dealCheck()
	print('lua: pos1 force is '..com.RealDt.CR[1].MaxForce)
	print('lua: pos2 force is '..com.RealDt.CR[2].MaxForce)
	print('lua: pos3 force is '..com.RealDt.CR[3].MaxForce)
end

function DealPush()
	local tstFlg = true
	local t1 = GetMilliseconds()
	local t2 = t1
	
	com.PlcApi_Release(false, TaskConfig.Pos[1])
	com.PlcApi_Menu_GainRate(true, TaskConfig.Pos[1], 2)
	com.PlcApi_Push_Target(false, TaskConfig.Pos[1],0,0, 60 , 20)
	while  tstFlg do 
		dealCheck()
		if com.RealDt.CR[2].MaxForce < 57 then
			com.PlcApi_Nothing(false)
		else
			tstFlg = false
		end
	end
	
	com.PlcApi_Release(true, TaskConfig.Pos[1])
	com.PlcApi_Over()
	print("over")
end

function DealPullPush()
	local tstCnt = TaskConfig.TotalTimes
	local t1 = GetMilliseconds()
	local t2 = t1
	
	while  tstCnt > 0 do 
		LogTaskInfo.Index = LogTaskInfo.Index+1
		com.PlcApi_Push_Fix(true, TaskConfig.Pos[1],0,0,1200)
		t1 = GetMilliseconds()
		t2 = t1
		while t2-t1 < 1000 do
			com.PlcApi_Nothing(false)
			t2 =  GetMilliseconds()
		end
		--print(json.stringify(TaskConfig))
		com.PlcApi_LogPoint(LogTaskInfo, LogTaskPosFields)
		
		
		com.PlcApi_Pull_Fix(true, TaskConfig.Pos[1],0,0,1200)
		t1 = GetMilliseconds()
		t2 = t1
		while t2-t1 < 1000 do
			com.PlcApi_Nothing(false)
			t2 =  GetMilliseconds()
		end
		--print(json.stringify(TaskConfig))
		com.PlcApi_LogPoint(LogTaskInfo, LogTaskPosFields)
		
		tstCnt = tstCnt -1
	end
	
	com.PlcApi_Release(true, TaskConfig.Pos[1])
	com.PlcApi_Over()
	print("over")
end


--下面这部分功能可能并不重要，用户可以根据停下来的情况，手动设置次数来继续
function StoreEnv()
	return Common.StoreEnv(TaskConfig,RunInfo)
end

function RestoreEnv(envJson)
	TaskConfig,RunInfo = Common.RestoreEnv(envJson)
end