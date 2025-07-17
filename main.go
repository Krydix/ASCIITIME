package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	"os/signal"
	"syscall"
	"golang.org/x/term"
	"regexp"
)

// Color codes for cyberpunk 80s aesthetic with triadic color theory
const (
	// Primary colors (triadic scheme: magenta, cyan, yellow)
	NEON_MAGENTA = "\033[95m"  // Bright magenta
	NEON_CYAN    = "\033[96m"  // Bright cyan
	NEON_YELLOW  = "\033[93m"  // Bright yellow
	
	// Secondary colors for depth
	ORANGE       = "\033[38;5;208m"  // Sunset orange
	PINK         = "\033[38;5;213m"  // Hot pink
	PURPLE       = "\033[38;5;141m"  // Deep purple
	
	// Accent colors
	ELECTRIC_BLUE = "\033[38;5;81m"   // Electric blue
	LIME_GREEN    = "\033[38;5;154m"  // Lime green
	HOT_PINK      = "\033[38;5;198m"  // Hot pink
	
	// Additional colors for themes
	MATRIX_GREEN  = "\033[38;5;46m"   // Matrix green
	AMBER         = "\033[38;5;214m"  // Amber/orange
	RED           = "\033[38;5;196m"  // Bright red
	WHITE         = "\033[97m"        // Bright white
	BLUE          = "\033[38;5;39m"   // Blue
	VIOLET        = "\033[38;5;129m"  // Violet
	GOLD          = "\033[38;5;220m"  // Gold
	
	// Background and structure
	DARK_GRAY     = "\033[38;5;236m"  // Dark gray
	RESET         = "\033[0m"         // Reset
	BOLD          = "\033[1m"         // Bold
	DIM           = "\033[2m"         // Dim
)

// Color themes
type ColorTheme struct {
	Name        string
	Colors      []string
	BorderColor string
	TitleColor  string
	DateColor   string
	FooterColor string
	ColonColor  string
}

var colorThemes = []ColorTheme{
	{
		Name:        "Synthwave",
		Colors:      []string{NEON_MAGENTA, NEON_CYAN, NEON_YELLOW, ORANGE, PINK, PURPLE},
		BorderColor: NEON_CYAN,
		TitleColor:  NEON_MAGENTA,
		DateColor:   NEON_YELLOW,
		FooterColor: ELECTRIC_BLUE,
		ColonColor:  HOT_PINK,
	},
	{
		Name:        "Matrix",
		Colors:      []string{MATRIX_GREEN, MATRIX_GREEN, MATRIX_GREEN},
		BorderColor: MATRIX_GREEN,
		TitleColor:  MATRIX_GREEN,
		DateColor:   MATRIX_GREEN,
		FooterColor: MATRIX_GREEN,
		ColonColor:  MATRIX_GREEN,
	},
	{
		Name:        "Amber Terminal",
		Colors:      []string{AMBER, AMBER, AMBER},
		BorderColor: AMBER,
		TitleColor:  AMBER,
		DateColor:   AMBER,
		FooterColor: AMBER,
		ColonColor:  AMBER,
	},
	{
		Name:        "Neon Blue",
		Colors:      []string{ELECTRIC_BLUE, NEON_CYAN, BLUE},
		BorderColor: NEON_CYAN,
		TitleColor:  ELECTRIC_BLUE,
		DateColor:   BLUE,
		FooterColor: NEON_CYAN,
		ColonColor:  ELECTRIC_BLUE,
	},
	{
		Name:        "Sunset",
		Colors:      []string{RED, ORANGE, NEON_YELLOW, PINK},
		BorderColor: ORANGE,
		TitleColor:  RED,
		DateColor:   NEON_YELLOW,
		FooterColor: PINK,
		ColonColor:  RED,
	},
	{
		Name:        "Monochrome",
		Colors:      []string{WHITE, WHITE, WHITE},
		BorderColor: WHITE,
		TitleColor:  WHITE,
		DateColor:   WHITE,
		FooterColor: WHITE,
		ColonColor:  WHITE,
	},
	{
		Name:        "Purple Rain",
		Colors:      []string{PURPLE, VIOLET, NEON_MAGENTA},
		BorderColor: PURPLE,
		TitleColor:  NEON_MAGENTA,
		DateColor:   VIOLET,
		FooterColor: PURPLE,
		ColonColor:  NEON_MAGENTA,
	},
	{
		Name:        "Crimson",
		Colors:      []string{RED, RED, RED},
		BorderColor: RED,
		TitleColor:  RED,
		DateColor:   RED,
		FooterColor: RED,
		ColonColor:  RED,
	},
	{
		Name:        "Electric Blue",
		Colors:      []string{ELECTRIC_BLUE, ELECTRIC_BLUE, ELECTRIC_BLUE},
		BorderColor: ELECTRIC_BLUE,
		TitleColor:  ELECTRIC_BLUE,
		DateColor:   ELECTRIC_BLUE,
		FooterColor: ELECTRIC_BLUE,
		ColonColor:  ELECTRIC_BLUE,
	},
	{
		Name:        "Lime",
		Colors:      []string{LIME_GREEN, LIME_GREEN, LIME_GREEN},
		BorderColor: LIME_GREEN,
		TitleColor:  LIME_GREEN,
		DateColor:   LIME_GREEN,
		FooterColor: LIME_GREEN,
		ColonColor:  LIME_GREEN,
	},
	{
		Name:        "Gold",
		Colors:      []string{GOLD, GOLD, GOLD},
		BorderColor: GOLD,
		TitleColor:  GOLD,
		DateColor:   GOLD,
		FooterColor: GOLD,
		ColonColor:  GOLD,
	},
	{
		Name:        "Pink Dream",
		Colors:      []string{HOT_PINK, HOT_PINK, HOT_PINK},
		BorderColor: HOT_PINK,
		TitleColor:  HOT_PINK,
		DateColor:   HOT_PINK,
		FooterColor: HOT_PINK,
		ColonColor:  HOT_PINK,
	},
	{
		Name:        "Violet Mono",
		Colors:      []string{VIOLET, VIOLET, VIOLET},
		BorderColor: VIOLET,
		TitleColor:  VIOLET,
		DateColor:   VIOLET,
		FooterColor: VIOLET,
		ColonColor:  VIOLET,
	},
	{
		Name:        "Cyan Dream",
		Colors:      []string{NEON_CYAN, NEON_CYAN, NEON_CYAN},
		BorderColor: NEON_CYAN,
		TitleColor:  NEON_CYAN,
		DateColor:   NEON_CYAN,
		FooterColor: NEON_CYAN,
		ColonColor:  NEON_CYAN,
	},
	{
		Name:        "Magenta Mono",
		Colors:      []string{NEON_MAGENTA, NEON_MAGENTA, NEON_MAGENTA},
		BorderColor: NEON_MAGENTA,
		TitleColor:  NEON_MAGENTA,
		DateColor:   NEON_MAGENTA,
		FooterColor: NEON_MAGENTA,
		ColonColor:  NEON_MAGENTA,
	},
}

// Enhanced ASCII art for digits 0-9 with better readability and cyberpunk aesthetics
var asciiDigits = map[rune][]string{
	'0': {
		"  ████████  ",
		"  ██    ██  ",
		"  ██    ██  ",
		"  ██    ██  ",
		"  ██    ██  ",
		"  ████████  ",
	},
	'1': {
		"      ██    ",
		"    ████    ",
		"      ██    ",
		"      ██    ",
		"      ██    ",
		"  ██████████",
	},
	'2': {
		"  ████████  ",
		"  ██    ██  ",
		"      ████  ",
		"    ████    ",
		"  ████      ",
		"  ██████████",
	},
	'3': {
		"  ████████  ",
		"        ██  ",
		"    ██████  ",
		"        ██  ",
		"        ██  ",
		"  ████████  ",
	},
	'4': {
		"  ██    ██  ",
		"  ██    ██  ",
		"  ██████████",
		"        ██  ",
		"        ██  ",
		"        ██  ",
	},
	'5': {
		"  ██████████",
		"  ██        ",
		"  ████████  ",
		"        ██  ",
		"  ██    ██  ",
		"  ████████  ",
	},
	'6': {
		"  ████████  ",
		"  ██        ",
		"  ████████  ",
		"  ██    ██  ",
		"  ██    ██  ",
		"  ████████  ",
	},
	'7': {
		"  ██████████",
		"        ██  ",
		"      ████  ",
		"    ████    ",
		"  ████      ",
		"  ██        ",
	},
	'8': {
		"  ████████  ",
		"  ██    ██  ",
		"  ████████  ",
		"  ██    ██  ",
		"  ██    ██  ",
		"  ████████  ",
	},
	'9': {
		"  ████████  ",
		"  ██    ██  ",
		"  ████████  ",
		"        ██  ",
		"        ██  ",
		"  ████████  ",
	},
	':': {
		"     ",
		"  ██ ",
		"     ",
		"     ",
		"  ██ ",
		"     ",
	},
	'.': {
		"     ",
		"     ",
		"     ",
		"     ",
		"  ██ ",
		"     ",
	},
}

// Transition effects for cyberpunk aesthetic
var transitionFrames = []string{
	"░", "▒", "▓", "█",
}

// Cyberpunk decorative elements
var cyberpunkBorder = []string{
	"▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄",
	"▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀",
}

type TimeDisplay struct {
	prevTime     string
	currentTime  string
	transition   bool
	step         int
	currentTheme int
	showDate     bool
	keyPressed   chan rune
	forceRefresh bool
	terminalState *term.State
}

func NewTimeDisplay() *TimeDisplay {
	return &TimeDisplay{
		currentTheme: 0,
		showDate:     false,
		keyPressed:   make(chan rune, 10),
		forceRefresh: false,
	}
}

func (td *TimeDisplay) getCurrentTheme() ColorTheme {
	return colorThemes[td.currentTheme%len(colorThemes)]
}

func (td *TimeDisplay) nextTheme() {
	td.currentTheme = (td.currentTheme + 1) % len(colorThemes)
	td.forceRefresh = true
}

func stripANSI(input string) string {
	// Remove ANSI escape sequences
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(input, "")
}

func (td *TimeDisplay) getTerminalSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		// Default fallback size
		return 80, 24
	}
	return width, height
}

func (td *TimeDisplay) clear() {
	// More reliable clear screen method
	fmt.Print("\033[H\033[2J")
}

func (td *TimeDisplay) getDigitColor(index int) string {
	theme := td.getCurrentTheme()
	return theme.Colors[index%len(theme.Colors)]
}

func (td *TimeDisplay) renderTime(timeStr string, isTransition bool, step int) {
	theme := td.getCurrentTheme()
	lines := make([]string, 6)
	
	// Build each line of the ASCII art
	for i := 0; i < 6; i++ {
		var lineBuilder strings.Builder
		charIndex := 0
		
		for _, char := range timeStr {
			if ascii, exists := asciiDigits[char]; exists {
				colorCode := td.getDigitColor(charIndex)
				
				if isTransition && char != ':' && char != '.' && step < len(transitionFrames) {
					// Apply transition effect
					transitionLine := strings.ReplaceAll(ascii[i], "█", transitionFrames[step])
					lineBuilder.WriteString(colorCode + transitionLine + RESET)
				} else {
					// Apply gradient effect to digits
					line := ascii[i]
					if char == ':' || char == '.' {
						line = strings.ReplaceAll(line, "██", theme.ColonColor+"██"+RESET)
					} else {
						line = strings.ReplaceAll(line, "██", colorCode+"██"+RESET)
					}
					lineBuilder.WriteString(line)
				}
				charIndex++
			}
		}
		lines[i] = lineBuilder.String()
	}
	
	// Simple header
	fmt.Print("\r\n")
	fmt.Printf("%s%s╔═══════════════════════════════════════════════════════════════════╗%s\r\n", theme.BorderColor, BOLD, RESET)
	
	var title string
	if td.showDate {
		title = "DATE-NEXUS"
	} else {
		title = "CHRONO-NEXUS"
	}
	
	fmt.Printf("%s%s║                           %s                           ║%s\r\n", theme.BorderColor, BOLD, title, RESET)
	fmt.Printf("%s%s╚═══════════════════════════════════════════════════════════════════╝%s\r\n", theme.BorderColor, BOLD, RESET)
	fmt.Print("\r\n")
	
	// Print ASCII art lines - no padding, just print them
	for _, line := range lines {
		fmt.Printf("      %s\r\n", line)
	}
	
	fmt.Print("\r\n")
	
	// Simple footer
	now := time.Now()
	var dateStr string
	if td.showDate {
		dateStr = now.Format("Monday, 02.01.2006 • 15:04:05")
	} else {
		dateStr = now.Format("Monday, 02.01.2006")
	}
	
	fmt.Printf("              %s%s%s%s\r\n", theme.DateColor, DIM, dateStr, RESET)
	fmt.Printf("           %s%sTheme: %s | Tab: Change | Space: Toggle%s\r\n", theme.FooterColor, DIM, theme.Name, RESET)
	fmt.Print("\r\n")
}

func (td *TimeDisplay) shouldTransition(prev, curr string) bool {
	if len(prev) != len(curr) {
		return false
	}
	
	// Transition on minute changes (positions 3 and 4 in HH:MM:SS format)
	if len(prev) >= 8 && len(curr) >= 8 {
		return prev[3] != curr[3] || prev[4] != curr[4]
	}
	return false
}

func (td *TimeDisplay) handleKeyInput() {
	go func() {
		for {
			// Read single character without Enter
			var b = make([]byte, 1)
			os.Stdin.Read(b)
			
			char := rune(b[0])
			
			switch char {
			case '\t': // Tab key
				select {
				case td.keyPressed <- 't':
				default:
				}
			case ' ': // Spacebar
				select {
				case td.keyPressed <- 'e':
				default:
				}
			case 't', 'T': // 't' key
				select {
				case td.keyPressed <- 't':
				default:
				}
			case 'e', 'E': // 'e' key
				select {
				case td.keyPressed <- 'e':
				default:
				}
			case 'q', 'Q', '\x03': // Ctrl+C or 'q'
				td.restoreTerminal()
				fmt.Print("\033[?25h")
				os.Exit(0)
			}
		}
	}()
}

func (td *TimeDisplay) setRawMode() {
	// Put terminal into raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("Error setting raw mode: %v\n", err)
		return
	}
	td.terminalState = oldState
	
	fmt.Print("\r\n" + NEON_CYAN + "Controls: Tab or 't': Change Theme | Spacebar or 'e': Toggle Date View | 'q' or Ctrl+C: Quit" + RESET + "\r\n")
	time.Sleep(1 * time.Second)
}

func (td *TimeDisplay) restoreTerminal() {
	if td.terminalState != nil {
		term.Restore(int(os.Stdin.Fd()), td.terminalState)
	}
}

func (td *TimeDisplay) update() {
	now := time.Now()
	var timeStr string
	
	if td.showDate {
		timeStr = now.Format("02.01.06")
	} else {
		timeStr = now.Format("15:04:05")
	}
	
	// Check if we need to transition (only for time view)
	if !td.showDate && td.currentTime != "" && td.shouldTransition(td.currentTime, timeStr) {
		td.prevTime = td.currentTime
		td.transition = true
		td.step = 0
	}
	
	td.currentTime = timeStr
	
	// Clear screen before rendering
	td.clear()
	
	// Handle transition animation (only for time view)
	if !td.showDate && td.transition {
		td.renderTime(td.currentTime, true, td.step)
		td.step++
		
		if td.step >= len(transitionFrames) {
			td.transition = false
			td.step = 0
		}
	} else if td.forceRefresh {
		// Force refresh on theme change or view toggle
		td.renderTime(td.currentTime, false, 0)
		td.forceRefresh = false
	} else {
		td.renderTime(td.currentTime, false, 0)
	}
}

func main() {
	// Hide cursor
	fmt.Print("\033[?25l")
	
	// Ensure cursor is shown when program exits
	defer fmt.Print("\033[?25h")
	
	display := NewTimeDisplay()
	
	// Set terminal to raw mode
	display.setRawMode()
	
	// Handle terminal restoration on exit
	defer display.restoreTerminal()
	
	// Set up signal handling for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		display.restoreTerminal()
		fmt.Print("\033[?25h")
		os.Exit(0)
	}()
	
	// Start keyboard input handling
	display.handleKeyInput()
	
	// Print startup message
	theme := display.getCurrentTheme()
	fmt.Printf("%s%s", theme.TitleColor, BOLD)
	fmt.Print("╔════════════════════════════════════════════════════════════════════╗\r\n")
	fmt.Print("║                    INITIALIZING CHRONO-NEXUS                       ║\r\n")
	fmt.Print("║                   TEMPORAL MATRIX LOADING...                       ║\r\n")
	fmt.Print("║                                                                    ║\r\n")
	fmt.Printf("║                   Current Theme: %-30s ║\r\n", theme.Name)
	fmt.Print("║                   Tab or 't': Change Theme                     ║\r\n")
	fmt.Print("║                   Space or 'e': Toggle Date View              ║\r\n")
	fmt.Print("║                   'q' or Ctrl+C: Quit                         ║\r\n")
	fmt.Print("╚════════════════════════════════════════════════════════════════════╝\r\n")
	fmt.Printf("%s", RESET)
	
	// Brief startup animation
	time.Sleep(2 * time.Second)
	
	// Create a ticker that updates every second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	
	// Initial display
	display.update()
	
	// Main loop
	for {
		select {
		case key := <-display.keyPressed:
			// Handle key press immediately
			switch key {
			case 't': // Tab key
				display.nextTheme()
				display.update()
			case 'e': // Enter key
				display.showDate = !display.showDate
				display.forceRefresh = true
				display.update()
			}
		case <-ticker.C:
			display.update()
		}
	}
}
