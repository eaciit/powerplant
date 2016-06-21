;WITH CTE_DB as 
(
    SELECT * FROM DataBrowser 
    WHERE 
        1=1 
        @PeriodYear 
        @EquipmentType 
),
CTE_PLANT as
(
    SELECT * FROM PowerPlantCoordinates 
    WHERE 
        1=1 
        @PlantName 
),
CTE_WO as 
(

    SELECT * from  
    (SELECT *, 
        (Select CAST((Select case 
                when DateDiff(MINUTE, WO.ScheduledStart,WO.ScheduledFinish) > 0  
                then DateDiff(MINUTE, WO.ScheduledStart,WO.ScheduledFinish) / 60.0 
                else 0 
            END) as float) 
        ) as PlanDuration, 
        (Select CAST((Select case 
                when DateDiff(MINUTE, WO.ActualStart, WO.ActualFinish) > 0 
                then DateDiff(MINUTE, WO.ActualStart, WO.ActualFinish) / 60.0 
                else 0 
            END) as float) 
        ) as ActualDuration 
    FROM WOList as WO 
    WHERE 
        1=1 
        and WO.ActualStart >= @PeriodFrom  
        and WO.ActualStart < @PeriodTo 
    ) as qr
    WHERE 1=1
)

SELECT 
    @Summary
    FROM
    (Select
        DataBrowser.Id as Id,
        DataBrowser.PeriodYear as PeriodYear,
        DataBrowser.FunctionalLocation as FunctionalLocation,
        DataBrowser.FLDescription as FLDescription,
        DataBrowser.IsTurbine as IsTurbine,
        DataBrowser.IsSystem as IsSystem,
        DataBrowser.TurbineParent as TurbineParent,
        DataBrowser.SystemParent SystemParent,
        DataBrowser.AssetType as AssetType,
        DataBrowser.EquipmentType EquipmentType,
        DataBrowser.EquipmentTypeDescription as EquipmentTypeDescription,
        DataBrowser.PlantCode as PlantCode,
        DataBrowser.TInfShortName as TInfShortName,
        DataBrowser.TInfManufacturer as TInfManufacturer,
        DataBrowser.TInfModel as TInfModel,
        DataBrowser.TInfUnitType as TInfUnitType,
        DataBrowser.TInfInstalledCapacity as TInfInstalledCapacity,
        DataBrowser.TInfOperationalCapacity as TInfOperationalCapacity,
        DataBrowser.TInfPrimaryFuel as TInfPrimaryFuel,      
        DataBrowser.TInfPrimaryFuel2 as TInfPrimaryFuel2,     
        DataBrowser.TInfBackupFuel as TInfBackupFuel,       
        DataBrowser.TInfHeatRate as TInfHeatRate,
        DataBrowser.TInfEfficiency as TInfEfficiency,       
        DataBrowser.TInfCommisioningDate as TInfCommisioningDate, 
        DataBrowser.TInfRetirementPlan as TInfRetirementPlan,   
        DataBrowser.TInfInstalledMWH as TInfInstalledMWH,     
        DataBrowser.TInfActualEnergyGeneration as TInfActualEnergyGeneration,
        DataBrowser.TInfActualFuelConsumption_GASMMSCF as TInfActualFuelConsumption_GASMMSCF,
        DataBrowser.TInfActualFuelConsumption_CrudeBarrel as TInfActualFuelConsumption_CrudeBarrel,
        DataBrowser.TInfActualFuelConsumption_HFOBarrel as TInfActualFuelConsumption_HFOBarrel,
        DataBrowser.TInfActualFuelConsumption_DieselBarrel as TInfActualFuelConsumption_DieselBarrel,
        DataBrowser.TInfCapacityFactor as TInfCapacityFactor,   
        DataBrowser.TInfUpdateEnergyGeneration as TInfUpdateEnergyGeneration,
        DataBrowser.TInfUpdateFuelConsumption as TInfUpdateFuelConsumption, 

        Plant.PlantCode as PlantPlantCode,
        Plant.PlantName as PlantPlantName,
        Plant.PlantType as PlantPlantType,
        Plant.Province as PlantProvince,
        Plant.Region as PlantRegion,
        Plant.City as PlantCity,
        Plant.FuelTypes_Crude as PlantFuelTypes_Crude,
        Plant.FuelTypes_Heavy as PlantFuelTypes_Heavy,
        Plant.FuelTypes_Diesel as PlantFuelTypes_Diesel,
        Plant.FuelTypes_Gas as PlantFuelTypes_Gas,
        Plant.GasTurbineUnit as PlantGasTurbineUnit,
        Plant.GasTurbineCapacity as PlantGasTurbineCapacity,
        Plant.SteamUnit as PlantSteamUnit,
        Plant.SteamCapacity as PlantSteamCapacity,
        Plant.DieselUnit as PlantDieselUnit,
        Plant.DieselCapacity as PlantDieselCapacity,
        Plant.CombinedCycleUnit as PlantCombinedCycleUnit,
        Plant.CombinedCycleCapacity as PlantCombinedCycleCapacity,
        Plant.Longitude as PlantLongitude,
        Plant.Latitude as PlantLatitude,

        WO.Type as WorkOrderType,
        WO.FunctionalLocation as MaintenanceFunctionalLocation,	
        WO.OrderCode as MaintenanceOrder,
        WO.Description as MaintenanceDescription,
        WO.ScheduledStart as PlanStart,
        WO.ScheduledFinish as PlanFinish,
        WO.ActualStart as ActualStart,
        WO.ActualFinish as ActualFinish,
        WO.ActualCost as MaintenanceCost,

        WO.PlanDuration as PlanDuration,
        WO.ActualDuration as ActualDuration,
        (Select case
                when WO.ActualStart IS NOT NULL
                then (SELECT CAST((DATEDIFF(
                        HOUR,
                        (Select top 1 wl.ActualFinish from WOList as wl where wl.FunctionalLocation = DataBrowser.FunctionalLocation and wl.ActualStart < WO.ActualStart order by wl.ActualStart desc),
                        WO.ActualStart) / 24.0) as float))
                else 0
            END) as MaintenanceInterval 						
                    
    from CTE_DB DataBrowser 
    LEFT JOIN CTE_Plant as Plant 
    ON DataBrowser.PlantCode = Plant.PlantCode 
    LEFT Join CTE_WO WO 
    ON DataBrowser.FunctionalLocation = WO.FunctionalLocation 
    WHERE 
        DataBrowser.FunctionalLocation = WO.FunctionalLocation and 
        DataBrowser.PlantCode = Plant.PlantCode
)as result