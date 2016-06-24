USE ecsec
GO

-- Drop The index first to make sure the index is not exist
-- put in here the drop index script
drop index IDX_DTBROWSER_MAIN on DataBrowser
drop index IDX_PPC_MAIN on PowerPlantCoordinates
drop index IDX_WOL_MAIN on WOList
drop index IDX_MAINTENANCECOST_MAIN




-- Create the index
-- put in here the create index script
create index IDX_DTBROWSER_MAIN on DataBrowser (PeriodYear Asc,FunctionalLocation Asc,TurbineParent Asc,SystemParent Asc,EquipmentType Asc,PlantCode Asc);
create index IDX_PPC_MAIN on PowerPlantCoordinates (PlantCode Asc, PlantName Asc);
create index IDX_WOL_MAIN on WOList (FunctionalLocation Asc, ScheduledStart Asc, ScheduledFinish Asc, UserStatus Asc, Type Asc, OrderCode Asc, NotificationCode Asc, Plant Asc, EquipmentCode Asc, SortField Asc, ActualStart Desc, ActualFinish Asc, ActualCost Asc);
create index IDX_MAINTENANCECOST_MAIN on MaintenanceCost (MaintenanceActivityType Asc, MaintenanceOrder Asc);
