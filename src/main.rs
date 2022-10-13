use pervasives::relative;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    if relative!("stack.yaml")?.exists() {
        println!("haskell project found");
        return Ok(());
    }

    if relative!("package.json")?.exists() && relative!("pnpm-lock.yaml")?.exists() {
        println!("node pnpm project found");
        return Ok(());
    }

    if relative!("Cargo.toml")?.exists() {
        println!("rust project found");
        return Ok(());
    }

    if relative!("go.mod")?.exists() {
        println!("go project found");
        return Ok(());
    }

    Ok(())
}
