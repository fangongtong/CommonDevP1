
description = [[[
	这里是文件说明
]]]

-- 每个lua文件都应该有CylinderConfig对象用于定义任务对应的缸
CylinderConfig = {
	Cylinder = {
		{Name="left", Pos=nil},
		{Name="right", Pos=nil},
	},
}

-- 每个lua文件都应该有TaskConfig对象用于定义任务信息
TaskConfig = {
	
}

function initCylinder(cyl)
	print(cyl.Name)
end

function InitCylinder()
	for k,v in pairs(CylinderConfig.Cylinder) do
		initCylinder(v)
	end
end

function StartWork()
	print("StartWork")
	
	for k,v in pairs(CylinderConfig.Cylinder) do
		print(v.Name.." work on pos:",v.Pos)
	end
	
	return 0
end

function StopWork()
	print("StopWork")
	
	
	return 0
end
