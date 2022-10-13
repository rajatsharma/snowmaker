#[macro_export]
macro_rules! relative {
    ($e:literal) => {
        ::std::env::current_dir().map(|path| path.join($e))
    };
}
