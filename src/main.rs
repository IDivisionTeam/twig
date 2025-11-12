mod cmd;
mod network;

use anyhow::Result;
use cmd::twig;

fn main() -> Result<()> {
    twig::execute()
}
