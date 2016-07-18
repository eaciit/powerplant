
#get actual duration

CREATE PROCEDURE spDataBrowserGetActualDuration
   @WOType nvarchar(10)
AS
BEGIN

select (Select CAST((Select case
   when WO.ActualFinish IS NOT NULL AND WO.ActualStart IS NOT NULL
      Then (Select case
         when DateDiff(SECOND, WO.ActualStart, WO.ActualFinish) > 0
            then DateDiff(SECOND, WO.ActualStart, WO.ActualFinish) / 3600.000000000000000
         else 0
      END)
   Else 0
   END) as float )
) as ActualDuration from DataBrowser as db inner join wolist as wo on db.FunctionalLocation = wo.FunctionalLocation where  wo.Type = @WOType

END