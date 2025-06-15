
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



function CheckCylinderConfig(cfg)
	for k,v in pairs(CylinderConfig.Cylinder) do
		if v.Pos == nil then 
			return false
		end
		--print(v.Name.." work on pos:",v.Pos)
	end
	return true
end

json = require "tasks.json"

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

function common.PlcApi_Nothing(needRtn)
	coroutine.yield(needRtn, 0, 1)
end

function common.PlcApi_Over()
	coroutine.yield(false, 0, 0)
end
return common
