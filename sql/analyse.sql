create table if not exists pagination_test_result (
    pagination_id bigint,
    pos integer,
    analyse_type varchar(36),
    primary key (pagination_id,pos,analyse_type),
    val decimal(15,6),
    constraint `fk_papervote_optical_test_result_pagination_id` foreign key (`pagination_id`) references `papervote_optical` (`pagination_id`) on delete cascade on update cascade

);

with bv as (
    select pagination_id,pos,analyse_type,val from pagination_test_result where analyse_type="bildverarbeitung"
), tf as (
    select pagination_id,pos,analyse_type,val from pagination_test_result where analyse_type="tensorflow 32x32x150"
), ln as (
    select pagination_id,pos,analyse_type,val from pagination_test_result where analyse_type="local nn"
), wz as (
    select pagination_id,result_index - 1 pos,"zählung" analyse_type, if(marked='O',0,if(marked='W',0.5,1)) val from view_papervote_optical_result_ballotpaper
)
select 
    bv.pagination_id,
    bv.pos,
    bv.val bv,
    tf.val tf,
    ln.val ln,
    wz.val wz
from 
    bv
    join tf on (bv.pagination_id,bv.pos) = (tf.pagination_id,tf.pos)
    join wz on (bv.pagination_id,bv.pos) = (wz.pagination_id,wz.pos)
    join ln on (bv.pagination_id,bv.pos) = (ln.pagination_id,ln.pos)
having bv=0 and tf > 1