
local common = {}


common.RealDt = {
	--Break = false,
	DevAlarm = 0,
	CR = {
		{
			Alarm = 0,
			RealForce = 0,
			MaxForce = 0,
			RealDisplace = 0,
			Threshold = 0
		},
		{
			Alarm = 0,
			RealForce = 0,
			MaxForce = 0,
			RealDisplace = 0,
			Threshold = 0
		},
		{
			Alarm = 0,
			RealForce = 0,
			MaxForce = 0,
			RealDisplace = 0,
			Threshold = 0
		},
		{
			Alarm = 0,
			RealForce = 0,
			MaxForce = 0,
			RealDisplace = 0,
			Threshold = 0
		},
		{
			Alarm = 0,
			RealForce = 0,
			MaxForce = 0,
			RealDisplace = 0,
			Threshold = 0
		},
		{
			Alarm = 0,
			RealForce = 0,
			MaxForce = 0,
			RealDisplace = 0,
			Threshold = 0
		},
		{
			Alarm = 0,
			RealForce = 0,
			MaxForce = 0,
			RealDisplace = 0,
			Threshold = 0
		},
		{
			Alarm = 0,
			RealForce = 0,
			MaxForce = 0,
			RealDisplace = 0,
			Threshold = 0
		},
	},
}

local json = require "tasks.json"
--local util = require "tasks.cjson.util"


function common.SetTaskConfig(jsStr)
	--local res = json.decode(jsStr)
	local res = json.parse(jsStr)
	print(json.stringify(res))
	for idx,val in pairs(com.TaskConfig.Pos) do 
		com.TaskConfig.Pos[idx] = res.Pos[idx]
	end
	for idx,val in pairs(com.TaskConfig.PosLimitSw) do 
		com.TaskConfig.PosLimitSw[idx] = res.PosLimitSw[idx]
	end
	com.SetNewTaskConfig(res)
end

function common.GetTaskUsingRes()
	--return json.encode(com.TaskConfig)
	return json.stringify(com.TaskConfig.Pos), json.stringify(com.TaskConfig.PosLimitSw)
end

function CheckCylinderConfig(cfg)
	for k,v in pairs(CylinderConfig.Cylinder) do
		if v.Pos == nil then 
			return false
		end
		--print(v.Name.." work on pos:",v.Pos)
	end
	return true
end

function StoreEnv(tskEnv,runInfo)
	return json.encode({tskEnv,runInfo})
end

function RestoreEnv(envJson)
	return json.decode(envJson)
end

--RunInfo = {
--	CommuTime = 20	--ms
--}

common.Err_CylinderConfig = 1
common.Err_TaskConfig = 2


function common.PlcApi_Release(needRtn, pos) 
	coroutine.yield(needRtn, 1, pos)
end

function common.PlcApi_Push_Fix(needRtn, pos, sw1, sw2, threshold) 
	coroutine.yield(needRtn, 3, pos,sw1, sw2, threshold)
end
function common.PlcApi_Pull_Fix(needRtn, pos, sw1, sw2, threshold) 
	coroutine.yield(needRtn, 4, pos,sw1, sw2, threshold)
end

function common.PlcApi_Push_Target(needRtn, pos, sw1, sw2, targetF, threshold) 
	coroutine.yield(needRtn, 5, pos,sw1, sw2, targetF, threshold)
end
function common.PlcApi_Pull_Target(needRtn, pos, sw1, sw2, targetF, threshold) 
	coroutine.yield(needRtn, 6, pos,sw1, sw2, targetF, threshold)
end

function common.PlcApi_Menu_GainRate(needRtn,pos, gainVal)
	coroutine.yield(needRtn, 20, pos, gainVal)
end

function common.PlcApi_Nothing(needRtn)
	coroutine.yield(needRtn, 0, 1)
end

function common.PlcApi_Over()
	coroutine.yield(false, 0, 0)
end
return common
