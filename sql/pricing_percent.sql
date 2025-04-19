select
avg(面料占比) 平均面料占比
,max(面料占比) 最大面料占比
,min(面料占比) 最小面料占比

,avg(辅料占比) 平均辅料占比
,max(辅料占比) 最大辅料占比
,min(辅料占比) 最小辅料占比

,avg(包装占比) 平均包装占比
,max(包装占比) 最大包装占比
,min(包装占比) 最小包装占比

,avg(外协占比) 平均外协占比
,max(外协占比) 最大外协占比
,min(外协占比) 最小外协占比

,avg(加工费占比) 平均加工费占比
,max(加工费占比) 最大加工费占比
,min(加工费占比) 最小加工费占比

,avg(利税占比) 平均利税占比
,max(利税占比) 最大利税占比
,min(利税占比) 最小利税占比

,avg(其他占比) 平均其他占比
,max(其他占比) 最大其他占比
,min(其他占比) 最小其他占比
from (  
select  to_char(calculate_date,'yyyy') 年份 
,case when cp.clothing_material_type_name in ('配件','游泳') then cp.clothing_material_type_name  else cp.needle_tatt_name end 针梭织 
,cp.clothing_material_type_detail_name        外贸体系服装小类名  
, (cp.fabric_subtotal/cp.tax_inclusive_price)::numeric(12,2) 面料占比
, (cp.auxiliary_material_subtotal/cp.tax_inclusive_price)::numeric(12,2) 辅料占比
, (cp.packaging_material_subtotal/cp.tax_inclusive_price)::numeric(12,2) 包装占比  
, (cp.printing_subtotal/cp.tax_inclusive_price)::numeric(12,2) 外协占比  
, ((cp.sewing_cost+cp.manage_cost+cp.tax_point_price+other_subtotal)/cp.tax_inclusive_price)::numeric(12,2) 加工费占比 
, ((cp.profit_price+cp.tax_point_price)/cp.tax_inclusive_price)::numeric(12,2)         利税占比  
, ((cp.freight_fee)/cp.tax_inclusive_price)::numeric(12,2)运费占比  
, ((cp.tax_inclusive_price-cp.fabric_subtotal -cp.auxiliary_material_subtotal -cp.packaging_material_subtotal -cp.printing_subtotal -(cp.sewing_cost+cp.manage_cost+cp.tax_point_price+other_subtotal)
        -(cp.profit_price+cp.tax_point_price) )/cp.tax_inclusive_price)::numeric(12,2)其他占比
FROM
        adm.adm_trd_sc_ord_check_price_life_cycle_nd cp 
where cp.cost_actual_finish_time is not null  
 and price_type_name='大货核价' AND cp.process_template_desc='已完成'
 and to_char(calculate_date,'yyyy')>='2023'
 and cp.order_sku_qty<=300
 ---%s---
 )tt