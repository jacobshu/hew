use ratatui::{
    layout::Alignment,
    prelude::{Constraint, Direction, Layout},
    style::{Color, Style},
    widgets::{Block, BorderType, Borders, Paragraph},
};

use crate::task::{app::App, tui::Frame};

pub fn render(app: &App, f: &mut Frame<'_>) {
    let chunks = Layout::default()
        .direction(Direction::Horizontal)
        .constraints(
            [
                Constraint::Percentage(50),
                Constraint::Percentage(50),
            ]
        .as_ref())
    .split(f.size());

    f.render_widget(
        Paragraph::new(format!(
            "
                Press `Esc`, `Ctrl-C` or `q` to stop running. \n\
                Press `j` and `k` to increment and decrement the counter.\n\
                Counter: {}
            ",
            app.counter
        ))
        .block(
            Block::default()
                .title("Counter App")
                .title_alignment(Alignment::Center)
                .borders(Borders::ALL)
                .border_type(BorderType::Rounded),
        )
        .style(Style::default().fg(Color::Blue))
        .alignment(Alignment::Center),
        chunks[0],
    );
}
