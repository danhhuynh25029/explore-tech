use utils;

fn main() {
    println!("{}",utils::add(1, 2));

    let a: i64 = 10;
    let b: i64 = 20;
    match utils::sub(&a,b) {
        Ok(v) => println!("{}",v),
        Err(e) => println!("error :{}",e)
    }

    utils::log::Log{ message: "1211313".parse().unwrap(), error_code: 1213 };
}
