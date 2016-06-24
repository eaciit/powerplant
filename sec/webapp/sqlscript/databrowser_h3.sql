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
        @PeriodYear 
        @EquipmentType 
        @FILTERS_DB
),
CTE_PLANT as
(
    SELECT * FROM( 
        SELECT 
        Plant.PlantCode as PlantPlantCode,
        Plant.PlantName as PlantPlantName 
        FROM PowerPlantCoordinates as Plant 
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
            END ) as float ) 
            ) as PlanDuration,
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
        FROM WOList as WO 
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
                        from WOList as wl 
                        where 
                            wl.FunctionalLocation = DataBrowser.FunctionalLocation and 
                            wl.ActualStart < WO.ActualStart 
                        order by wl.ActualStart desc 
                    ) IS NOT NULL
            then (SELECT CAST((DATEDIFF(
                    SECOND,
                    (Select top 1 wl.ActualFinish from WOList as wl where wl.FunctionalLocation = DataBrowser.FunctionalLocation and wl.ActualStart < WO.ActualStart order by wl.ActualStart desc),
                    WO.ActualStart) / 86400.000000000000000) as float ))
            else 0
        END) as MaintenanceInterval
        
        from CTE_DB DataBrowser, CTE_Plant as Plant, CTE_WO as WO 
        WHERE 
            DataBrowser.FunctionalLocation = WO.MaintenanceFunctionalLocation  
            and DataBrowser.PlantCode = Plant.PlantPlantCode 
)

SELECT * 
    FROM CTE_MAIN 
        @FILTERS_MAIN 
    ORDER BY @ORDERBY 
        OFFSET @Offset ROWS 
        FETCH NEXT @Limit ROWS ONLY