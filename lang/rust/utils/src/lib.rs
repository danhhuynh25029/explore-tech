pub mod log;

pub fn add(left: u64, right: u64) -> u64 {
    left + right
}

pub fn sub(left: &i64, right: i64) -> Result<i64, &str> {
    let r: i64 = *left - right;
    // let r: u64 = 9;
    if r <= 0 {
        return Err("not valid");
    }
    Ok(r)
}

#[cfg(test)]
mod tests {
    use super::*;
    #[test]
    fn it_works() {
        let result = add(2, 2);
        assert_eq!(result, 4);
    }
}
