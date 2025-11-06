package diff

import (
	"fmt"
	"html"
	"strings"
)

// Renderer is the interface for diff renderers
type Renderer interface {
	Render(result *DiffResult) (string, error)
}

// HTMLRenderer renders diff as HTML
type HTMLRenderer struct {
	ShowStats bool
}

// NewHTMLRenderer creates a new HTML renderer
func NewHTMLRenderer(showStats bool) *HTMLRenderer {
	return &HTMLRenderer{ShowStats: showStats}
}

// Render renders the diff result as HTML
func (r *HTMLRenderer) Render(result *DiffResult) (string, error) {
	var sb strings.Builder

	// HTML header
	sb.WriteString(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Document Diff</title>
    <style>
        body { font-family: 'Segoe UI', Arial, sans-serif; margin: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; border-bottom: 3px solid #4CAF50; padding-bottom: 10px; }
        .stats { background: #f9f9f9; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .stats-item { display: inline-block; margin-right: 30px; }
        .stats-label { font-weight: bold; color: #666; }
        .stats-value { color: #333; font-size: 1.2em; }
        .diff-line { padding: 8px 12px; margin: 2px 0; font-family: 'Courier New', monospace; border-left: 4px solid transparent; }
        .added { background-color: #e6ffed; border-left-color: #28a745; }
        .deleted { background-color: #ffeef0; border-left-color: #dc3545; text-decoration: line-through; }
        .modified { background-color: #fff3cd; border-left-color: #ffc107; }
        .unchanged { color: #666; }
        .position { color: #999; font-size: 0.9em; margin-right: 10px; }
        .legend { margin: 20px 0; padding: 10px; background: #f0f0f0; border-radius: 5px; }
        .legend-item { display: inline-block; margin-right: 20px; }
        .legend-color { display: inline-block; width: 20px; height: 20px; margin-right: 5px; vertical-align: middle; border-radius: 3px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Document Comparison</h1>
        <p><strong>Old:</strong> ` + html.EscapeString(result.OldDocument) + `</p>
        <p><strong>New:</strong> ` + html.EscapeString(result.NewDocument) + `</p>
`)

	// Stats section
	if r.ShowStats {
		sb.WriteString(`
        <div class="stats">
            <h2>Statistics</h2>
            <div class="stats-item">
                <span class="stats-label">Total Changes:</span>
                <span class="stats-value">` + fmt.Sprintf("%d", result.Stats.TotalChanges) + `</span>
            </div>
            <div class="stats-item">
                <span class="stats-label">Added:</span>
                <span class="stats-value" style="color: #28a745;">` + fmt.Sprintf("%d", result.Stats.AddedLines) + `</span>
            </div>
            <div class="stats-item">
                <span class="stats-label">Deleted:</span>
                <span class="stats-value" style="color: #dc3545;">` + fmt.Sprintf("%d", result.Stats.DeletedLines) + `</span>
            </div>
            <div class="stats-item">
                <span class="stats-label">Modified:</span>
                <span class="stats-value" style="color: #ffc107;">` + fmt.Sprintf("%d", result.Stats.ModifiedLines) + `</span>
            </div>
        </div>
`)
	}

	// Legend
	sb.WriteString(`
        <div class="legend">
            <strong>Legend:</strong>
            <span class="legend-item"><span class="legend-color" style="background: #e6ffed;"></span>Added</span>
            <span class="legend-item"><span class="legend-color" style="background: #ffeef0;"></span>Deleted</span>
            <span class="legend-item"><span class="legend-color" style="background: #fff3cd;"></span>Modified</span>
        </div>
`)

	// Changes section
	sb.WriteString(`
        <h2>Changes</h2>
        <div class="diff">
`)

	if len(result.Changes) == 0 {
		sb.WriteString(`<p style="color: #28a745; font-weight: bold;">✓ No changes detected - documents are identical</p>`)
	} else {
		for _, change := range result.Changes {
			sb.WriteString(r.renderChange(change))
		}
	}

	// HTML footer
	sb.WriteString(`
        </div>
    </div>
</body>
</html>
`)

	return sb.String(), nil
}

// renderChange renders a single change as HTML
func (r *HTMLRenderer) renderChange(change Change) string {
	var class string
	var text string

	switch change.Type {
	case DiffAdded:
		class = "added"
		text = html.EscapeString(change.New)
	case DiffDeleted:
		class = "deleted"
		text = html.EscapeString(change.Old)
	case DiffModified:
		class = "modified"
		text = html.EscapeString(change.Old) + " → " + html.EscapeString(change.New)
	default:
		class = "unchanged"
		text = html.EscapeString(change.Old)
	}

	return fmt.Sprintf(`<div class="diff-line %s"><span class="position">Line %d:</span>%s</div>`,
		class, change.Position+1, text)
}

// MarkdownRenderer renders diff as Markdown
type MarkdownRenderer struct {
	ShowStats bool
}

// NewMarkdownRenderer creates a new Markdown renderer
func NewMarkdownRenderer(showStats bool) *MarkdownRenderer {
	return &MarkdownRenderer{ShowStats: showStats}
}

// Render renders the diff result as Markdown
func (r *MarkdownRenderer) Render(result *DiffResult) (string, error) {
	var sb strings.Builder

	// Header
	sb.WriteString("# Document Comparison\n\n")
	sb.WriteString(fmt.Sprintf("**Old:** %s  \n", result.OldDocument))
	sb.WriteString(fmt.Sprintf("**New:** %s\n\n", result.NewDocument))

	// Stats
	if r.ShowStats {
		sb.WriteString("## Statistics\n\n")
		sb.WriteString(fmt.Sprintf("- **Total Changes:** %d\n", result.Stats.TotalChanges))
		sb.WriteString(fmt.Sprintf("- **Added:** %d\n", result.Stats.AddedLines))
		sb.WriteString(fmt.Sprintf("- **Deleted:** %d\n", result.Stats.DeletedLines))
		sb.WriteString(fmt.Sprintf("- **Modified:** %d\n\n", result.Stats.ModifiedLines))
	}

	// Changes
	sb.WriteString("## Changes\n\n")

	if len(result.Changes) == 0 {
		sb.WriteString("✓ No changes detected - documents are identical\n")
	} else {
		for _, change := range result.Changes {
			sb.WriteString(r.renderChange(change))
		}
	}

	return sb.String(), nil
}

// renderChange renders a single change as Markdown
func (r *MarkdownRenderer) renderChange(change Change) string {
	switch change.Type {
	case DiffAdded:
		return fmt.Sprintf("**Line %d** `+` %s\n\n", change.Position+1, change.New)
	case DiffDeleted:
		return fmt.Sprintf("**Line %d** `-` ~~%s~~\n\n", change.Position+1, change.Old)
	case DiffModified:
		return fmt.Sprintf("**Line %d** `~` ~~%s~~ → %s\n\n", change.Position+1, change.Old, change.New)
	default:
		return ""
	}
}

// PlainTextRenderer renders diff as plain text
type PlainTextRenderer struct {
	ShowStats   bool
	ColorOutput bool
}

// NewPlainTextRenderer creates a new plain text renderer
func NewPlainTextRenderer(showStats, colorOutput bool) *PlainTextRenderer {
	return &PlainTextRenderer{
		ShowStats:   showStats,
		ColorOutput: colorOutput,
	}
}

// Render renders the diff result as plain text
func (r *PlainTextRenderer) Render(result *DiffResult) (string, error) {
	var sb strings.Builder

	// Header
	sb.WriteString("Document Comparison\n")
	sb.WriteString(strings.Repeat("=", 50) + "\n\n")
	sb.WriteString(fmt.Sprintf("Old: %s\n", result.OldDocument))
	sb.WriteString(fmt.Sprintf("New: %s\n\n", result.NewDocument))

	// Stats
	if r.ShowStats {
		sb.WriteString("Statistics:\n")
		sb.WriteString(fmt.Sprintf("  Total Changes: %d\n", result.Stats.TotalChanges))
		sb.WriteString(fmt.Sprintf("  Added:         %d\n", result.Stats.AddedLines))
		sb.WriteString(fmt.Sprintf("  Deleted:       %d\n", result.Stats.DeletedLines))
		sb.WriteString(fmt.Sprintf("  Modified:      %d\n\n", result.Stats.ModifiedLines))
	}

	// Changes
	sb.WriteString("Changes:\n")
	sb.WriteString(strings.Repeat("-", 50) + "\n\n")

	if len(result.Changes) == 0 {
		sb.WriteString("✓ No changes detected - documents are identical\n")
	} else {
		for _, change := range result.Changes {
			sb.WriteString(r.renderChange(change))
		}
	}

	return sb.String(), nil
}

// renderChange renders a single change as plain text
func (r *PlainTextRenderer) renderChange(change Change) string {
	prefix := ""
	symbol := " "

	switch change.Type {
	case DiffAdded:
		symbol = "+"
		prefix = "ADDED"
	case DiffDeleted:
		symbol = "-"
		prefix = "DELETED"
	case DiffModified:
		symbol = "~"
		prefix = "MODIFIED"
	}

	if change.Type == DiffModified {
		return fmt.Sprintf("[%s] Line %d: %s → %s\n", prefix, change.Position+1, change.Old, change.New)
	} else if change.Type == DiffAdded {
		return fmt.Sprintf("[%s] Line %d: %s %s\n", prefix, change.Position+1, symbol, change.New)
	} else if change.Type == DiffDeleted {
		return fmt.Sprintf("[%s] Line %d: %s %s\n", prefix, change.Position+1, symbol, change.Old)
	}

	return ""
}
