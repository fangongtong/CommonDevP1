
description = [[
	这里是文件说明
	本文件用于单缸测试
]]

-- 每个lua文件都应该有CylinderConfig对象用于定义任务对应的缸
CylinderConfig = {
	Cylinder = {
		{Name="cylinder", Pos=nil},
	},
}

-- 每个lua文件都应该有TaskConfig对象用于定义任务信息
TaskConfig = {
	TaskType = 1,	--1:pull 2:push
	Force = 200,
	Frequency = 2,
	TotalTimes = 40,
}

RunInfo = {
	CurrentTimes = 0,
	ForceDuration = 0,	--ms
	IdleDuration = 0, --ms
	RemainDuration = 0, --s
}

print(package.path)

Common = require "common" 

function StopWork()
	print("StopWork")
	
	return 0
end


function CheckConfig()
	if Common.CheckCylinderConfig(CylinderConfig) == false then
		return Common.Err_CylinderConfig
	end
	
	return 0
end

function InitCylinder()
	RunInfo.IdleDuration = ((1000 / TaskConfig.Frequency) - Common.RunInfo.CommuTime*2)/2
	RunInfo.ForceDuration = RunInfo.IdleDuration+Common.RunInfo.CommuTime
	print(RunInfo.IdleDuration,RunInfo.ForceDuration)
end


local cyliderRecordDt = {
	lastCount = 0,
	--  下面这个可能要换成数组
	lastMaxForce = 0,
	
	syncCount = 0,
	--  下面这个可能要换成数组
	syncMaxForce = 0,
	syncCrntForce = 0
}

function StartWork()
	print("StartWork")
	for key, value in pairs(RunInfo) do
		print('\t', key, value)
	end
	
	
	cyliderRecordDt.lastCount = 0
	cyliderRecordDt.lastMaxForce = 0
	cyliderRecordDt.syncCount = 0
	cyliderRecordDt.syncMaxForce = 0
	
	while  RunInfo.CurrentTimes < TaskConfig.TotalTimes do 
		RunInfo.CurrentTimes= RunInfo.CurrentTimes+1
		
		cyliderRecordDt.syncCount = RunInfo.CurrentTimes
		
		-- call go code 发力
		print("-->pull to plc")
		GoSleep(RunInfo.ForceDuration)
		-- call go code 卸力
		print("-->release to plc")
		GoSleep(RunInfo.IdleDuration)
	end
	
	return 0
end


function DealReal()
	print("Start deal real")
	
	--local realDt = {}
	
    while true do
		--_,_,v3 = coroutine.yield()
		print("DealReal 1")
		coroutine.yield()
		--Common.UnpackRealDt({coroutine.yield(cyliderRecordDt)},realDt)
		--print("DealReal 2")
		--print(realDt.Com2)
		if Common.RealDt.Break == true then 
			break
		end
		
		--  如果还没有更新次数,那就更新本次最大力值
		cyliderRecordDt.syncCrntForce = Common.RealDt.CR[1].RealForce
		if cyliderRecordDt.syncCount - cyliderRecordDt.lastCount < 2 then
			if cyliderRecordDt.syncMaxForce < Common.RealDt.CR[1].RealForce then
				cyliderRecordDt.syncMaxForce = Common.RealDt.CR[1].RealForce
			end
		else
			--  如果已经更新次数了,那就更新上次最大力值
			cyliderRecordDt.lastCount = cyliderRecordDt.syncCount -1
			cyliderRecordDt.lastMaxForce = cyliderRecordDt.syncMaxForce
			
			--print("lua.DealReal->"..#realDt.CR)
			cyliderRecordDt.syncMaxForce = Common.RealDt.CR[1].RealForce
			
				--print("h1")
				print("sync count:"..cyliderRecordDt.syncCount)
				print("sync last:"..cyliderRecordDt.lastCount.." and lastMax:"..cyliderRecordDt.lastMaxForce)
			
			
			for key, value in pairs(cyliderRecordDt) do
		        print('\t', key, value)
		    end
			
--			for k,v in pairs(cyliderRecordDt) do 
--					print("h1.2")
--				if type(v) == "table" then
--					print("h2")
--					print(unpack(v))
--				else
--					print("h3")
--					print(k,v)
--				end
--			end
		end
		
	end
end

--下面这部分功能可能并不重要，用户可以根据停下来的情况，手动设置次数来继续
function StoreEnv()
	return Common.StoreEnv(TaskConfig,RunInfo)
end

function RestoreEnv(envJson)
	TaskConfig,RunInfo = Common.RestoreEnv(envJson)
end