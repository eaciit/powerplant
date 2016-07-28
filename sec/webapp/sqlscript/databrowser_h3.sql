;WITH CTE_DB as 
(
    SELECT 
        MasterFunctionalLocation.FunctionalLocationCode as FunctionalLocation,
        MasterFunctionalLocation.Description as FLDescription,
        MasterFunctionalLocation.EquipmentType EquipmentType,
        MasterEquipmentType.Description as EquipmentTypeDescription,
        MasterFunctionalLocation.Plant as PlantCode 
    FROM MasterFunctionalLocation left join MasterEquipmentType on MasterFunctionalLocation.EquipmentType = MasterEquipmentType.Type 
    WHERE  
        @EquipmentType 
        @FILTERS_DB 
),
CTE_PLANT as
(
    SELECT * FROM( 
        SELECT 
        Plant.PlantCode as PlantPlantCode,
        Plant.PlantName as PlantPlantName 
        FROM MasterPowerPlant as Plant 
    ) res 
    @PlantName 
),
CTE_WO as 
(

    SELECT * from  
    (   SELECT 
            WO.Type as WorkOrderType,
            WO.FunctionalLocation as MaintenanceFunctionalLocation,
            WO.OrderCode as MaintenanceOrder,
            WO.Description as MaintenanceDescription,
            WO.ActualStart as ActualStart,
            WO.ActualFinish as ActualFinish,
            WO.ActualCost as MaintenanceCost, 
            (Select CAST((Select case 
                            when WO.ActualFinish IS NOT NULL AND WO.ActualStart IS NOT NULL 
                            Then (Select case 
                                    when DateDiff(SECOND, WO.ActualStart, WO.ActualFinish) > 0 
                                    then DateDiff(SECOND, WO.ActualStart, WO.ActualFinish) / 3600.000000000000000
                                    else 0 
                                END) 
                            Else 0 
            END) as float ) 
            ) as ActualDuration 
        FROM WorkOrderList as WO 
        WHERE 
            (select cast(WO.ActualStart as datetime)) >= (select cast(@PeriodFrom as datetime)) 
            and 
            (select cast(WO.ActualStart as datetime)) < (select cast(@PeriodTo as datetime)) 
    ) as RESULT 
    @FILTERS_WO 
), CTE_MAIN as (
    SELECT
        *,
        (SELECT case
            when WO.ActualStart IS NOT NULL AND (Select top 1 wl.ActualFinish 	
                        from WorkOrderList as wl 
                        where 
                            wl.FunctionalLocation = MasterFunctionalLocation.FunctionalLocation and 
                            wl.ActualStart < WO.ActualStart 
                        order by wl.ActualStart desc 
                    ) IS NOT NULL
            then (SELECT CAST((DATEDIFF(
                    SECOND,
                    (Select top 1 wl.ActualFinish from WorkOrderList as wl where wl.FunctionalLocation = MasterFunctionalLocation.FunctionalLocation and wl.ActualStart < WO.ActualStart order by wl.ActualStart desc),
                    WO.ActualStart) / 86400.000000000000000) as float ))
            else 0
        END) as MaintenanceInterval
        
        from CTE_DB MasterFunctionalLocation, CTE_Plant as Plant, CTE_WO as WO 
        WHERE 
            MasterFunctionalLocation.FunctionalLocation = WO.MaintenanceFunctionalLocation  
            and MasterFunctionalLocation.PlantCode = Plant.PlantPlantCode 
)

SELECT * 
    FROM CTE_MAIN 
        @FILTERS_MAIN 
    ORDER BY @ORDERBY 
        OFFSET @Offset ROWS 
        FETCH NEXT @Limit ROWS ONLY