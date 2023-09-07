use anyhow::Result;
use crossterm::{
    event::{self, Event::Key, KeyCode::Char},
    execute,
    terminal::{disable_raw_mode, enable_raw_mode, EnterAlternateScreen, LeaveAlternateScreen},
};
use ratatui::{
    prelude::{CrosstermBackend, Terminal},
    widgets::Paragraph,
};

pub type Frame<'a> = ratatui::Frame<'a, CrosstermBackend<std::io::Stderr>>;

struct App {
    counter: i64,
    should_quit: bool,
}

fn startup() -> Result<()> {
    enable_raw_mode()?;
    execute!(std::io::stderr(), crossterm::terminal::EnterAlternateScreen)?;
    Ok(())
}

fn update(app: &mut App) -> Result<()> {
    if event::poll(std::time::Duration::from_millis(250))? {
        if let event::Event::Key(key) = crossterm::event::read()? {
            match key.code {
                event::KeyCode::Char('j') => app.counter += 1,
                event::KeyCode::Char('k') => app.counter -= 1,
                event::KeyCode::Char('q') => app.should_quit = true,
                _ => (),
            }
        }
    }
    Ok(())
}

fn ui(f: &mut Frame<'_>, app: &App) {
    f.render_widget(
        Paragraph::new(format!("Counter: {}", app.counter)),
        f.size(),
    );
}

fn shutdown() -> Result<()> {
    execute!(std::io::stderr(), crossterm::terminal::LeaveAlternateScreen)?;
    disable_raw_mode()?;
    Ok(())
}

fn run() -> Result<()> {
    let mut t = Terminal::new(CrosstermBackend::new(std::io::stderr()))?;
    let mut app = App {
        counter: 0,
        should_quit: false,
    };

    loop {
        t.draw(|f| {
            ui(f, &app);
        })?;

        update(&mut app)?;

        if app.should_quit {
            break;
        }
    }
    Ok(())
}

pub fn init() -> Result<(), Box<dyn std::error::Error>> {
    startup()?;
    let status = run();
    shutdown()?;
    status?;
    Ok(())
}
