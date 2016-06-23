;WITH CTE_DB as 
(
    SELECT
        DataBrowser.Id as Id,
        DataBrowser.PeriodYear as PeriodYear,
        DataBrowser.FunctionalLocation as FunctionalLocation,
        DataBrowser.FLDescription as FLDescription,
        DataBrowser.EquipmentType EquipmentType,
        DataBrowser.EquipmentTypeDescription as EquipmentTypeDescription,
        DataBrowser.PlantCode as PlantCode
    FROM DataBrowser 
    WHERE 
        1=1 
        @PeriodYear 
        @EquipmentType 
),
CTE_PLANT as
(
    SELECT
        Plant.PlantCode as PlantPlantCode,
        Plant.PlantName as PlantPlantName
    FROM PowerPlantCoordinates as Plant
    WHERE 
        1=1 
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
            WO.ScheduledStart as PlanStart,
            WO.ScheduledFinish as PlanFinish,
            WO.ActualStart as ActualStart,
            WO.ActualFinish as ActualFinish,
            WO.ActualCost as MaintenanceCost, 
            (Select CAST((Select case 
                            when WO.ScheduledFinish IS NOT NULL AND WO.ScheduledStart IS NOT NULL 
                            Then (Select case 
                                    when DateDiff(SECOND, WO.ScheduledStart,WO.ScheduledFinish) > 0  
                                    then DateDiff(SECOND, WO.ScheduledStart,WO.ScheduledFinish) / 3600.000000000000000 
                                    else 0 
                                END)
                            Else 0 
            END ) as float) 
            ) as PlanDuration,
            (Select CAST((Select case 
                            when WO.ActualFinish IS NOT NULL AND WO.ActualStart IS NOT NULL 
                            Then (Select case 
                                    when DateDiff(SECOND, WO.ActualStart, WO.ActualFinish) > 0 
                                    then DateDiff(SECOND, WO.ActualStart, WO.ActualFinish) / 3600.000000000000000
                                    else 0 
                                END) 
                            Else 0 
            END) as float) 
            ) as ActualDuration  
        FROM WOList as WO 
        WHERE 
            (select cast(WO.ActualStart as datetime)) >= (select cast(@PeriodFrom as datetime)) 
            and 
            (select cast(WO.ActualStart as datetime)) < (select cast(@PeriodTo as datetime)) 
    ) as qr
    WHERE 1=1
)

SELECT 
    @Summary
    FROM
    (SELECT *,
        (SELECT case 
            When (res.LastMaintenanceDate IS NOT NULL and res.ActualStart IS NOT NULL)
            then (SELECT CAST((DATEDIFF(SECOND,res.LastMaintenanceDate,res.ActualStart) / 86400.000000000000000) as float))
            else 0 
			END
        ) as MaintenanceInterval 
        FROM
        (
            Select *,
                (Select top 1 wl.ActualFinish 	
					from WOList as wl 
					where 
						wl.FunctionalLocation = DataBrowser.FunctionalLocation and 
						wl.ActualStart < WO.ActualStart 
					order by wl.ActualStart desc 
				) as LastMaintenanceDate 
            from CTE_DB DataBrowser, CTE_Plant as Plant, CTE_WO as WO 
            WHERE 
                DataBrowser.FunctionalLocation = WO.MaintenanceFunctionalLocation and 
                DataBrowser.PlantCode = Plant.PlantPlantCode 
            
        )as res
    ) as result