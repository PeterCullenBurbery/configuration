package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "formatted_time_zone_info")
	if err != nil {
		log.Fatalf("❌ Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	javaFileName := "formatted_time_zone_info.java"
	className := "formatted_time_zone_info"
	javaFilePath := filepath.Join(tempDir, javaFileName)

	// Java source code
	javaCode := `import java.time.*;
import java.time.format.DateTimeFormatter;
import java.time.temporal.WeekFields;

public class formatted_time_zone_info {
    public static void main(String[] args) {
        ZonedDateTime now = ZonedDateTime.now();
        ZoneId tz = now.getZone();

        String date_part = now.format(DateTimeFormatter.ofPattern("yyyy-0MM-0dd"));
        String time_part = now.format(DateTimeFormatter.ofPattern("0HH.0mm.0ss.nnnnnnn"));

        WeekFields wf = WeekFields.ISO;
        int week = now.get(wf.weekOfWeekBasedYear());
        int weekday = now.get(wf.dayOfWeek());
        int iso_year = now.get(wf.weekBasedYear());
        int day_of_year = now.getDayOfYear();

        String output = String.format(
            "%s %s %04d-W%03d-%03d %04d-%03d",
            date_part, time_part, iso_year, week, weekday, now.getYear(), day_of_year
        );
        output = output.replace(time_part, time_part + " " + tz);
        System.out.println(output);
    }
}`

	// Write Java file
	err = os.WriteFile(javaFilePath, []byte(javaCode), 0644)
	if err != nil {
		log.Fatalf("❌ Failed to write Java file: %v", err)
	}

	// Compile Java file
	cmdCompile := exec.Command("javac", javaFileName)
	cmdCompile.Dir = tempDir
	if err := cmdCompile.Run(); err != nil {
		log.Fatalf("❌ Failed to compile Java file: %v", err)
	}

	// Run compiled class and capture output
	cmdRun := exec.Command("java", className)
	cmdRun.Dir = tempDir

	var out bytes.Buffer
	cmdRun.Stdout = &out
	cmdRun.Stderr = &out

	if err := cmdRun.Run(); err != nil {
		log.Fatalf("❌ Failed to run Java class: %v\nOutput:\n%s", err, out.String())
	}

	// Print only the timestamp output (no extra newline)
	_, err = os.Stdout.Write(out.Bytes())
	if err != nil {
		log.Fatalf("❌ Failed to write output: %v", err)
	}
}
