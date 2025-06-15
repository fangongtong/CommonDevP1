
local common = {}


function common.CheckCylinderConfig(cfg)
	for k,v in pairs(CylinderConfig.Cylinder) do
		if v.Pos == nil then 
			return false
		end
		--print(v.Name.." work on pos:",v.Pos)
	end
	return true
end

json = require "json"

function common.StoreEnv(tskEnv,runInfo)
	return json.encode({tskEnv,runInfo})
end

function common.RestoreEnv(envJson)
	return json.decode(envJson)
end

common.RunInfo = {
	CommuTime = 20,	--ms
}


common.RealDt = {
	Break = false,
	Com1 = 0,
	Com2 = 0,
	Com3 = 0,
	CR = {
		{RealForce=0,DriveDegree=0},
		{RealForce=0,DriveDegree=0},
		{RealForce=0,DriveDegree=0},
		{RealForce=0,DriveDegree=0},
		{RealForce=0,DriveDegree=0},
		{RealForce=0,DriveDegree=0},
		{RealForce=0,DriveDegree=0},
		{RealForce=0,DriveDegree=0},
		{RealForce=0,DriveDegree=0},
		{RealForce=0,DriveDegree=0},
		{RealForce=0,DriveDegree=0},
		{RealForce=0,DriveDegree=0},
	},
}

function common.UnpackRealDt(src, trg)

		print("UnpackRealDt 1")
	trg.Break = src[1]
	trg.Com1 = src[2]
	trg.Com2 = src[3]
	trg.Com3 = src[4]
	--print(dt[1],dt[2],dt[3])
	--print(unpack(dt))
	
	--print(unpack(trg))
	trg.CR = trg.CR or {}
	
	local j = 1
	for i=5,#src,2 do 
	
		trg.CR[j] = trg.CR[j] or {}
		trg.CR[j].RealForce = src[i]
		trg.CR[j].DriveDegree = src[i+1]
		j = j+1
	end 
	--for k,v in pairs(trg) do 
	--	print(k,v)
	--end 
end


common.Err_CylinderConfig = 1
common.Err_TaskConfig = 2

return common
