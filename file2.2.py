#!/usr/bin/env python3
"""
Enhanced data analysis utility for processing CSV and JSON files
"""

import csv
import json
import sys
from pathlib import Path
from typing import List, Dict, Optional


class DataProcessor:
    """Handles CSV and JSON data processing and analysis"""

    def __init__(self, filename: str):
        self.filename = filename
        self.data = []
        self.file_type = self._detect_file_type()

    def _detect_file_type(self) -> str:
        """Detect file type from extension"""
        suffix = Path(self.filename).suffix.lower()
        return 'json' if suffix == '.json' else 'csv'

    def load_data(self) -> None:
        """Load data from CSV or JSON file"""
        try:
            if self.file_type == 'json':
                with open(self.filename, 'r') as f:
                    self.data = json.load(f)
            else:
                with open(self.filename, 'r') as f:
                    reader = csv.DictReader(f)
                    self.data = list(reader)
        except FileNotFoundError:
            print(f"Error: File '{self.filename}' not found")
            sys.exit(1)
        except json.JSONDecodeError as e:
            print(f"Error: Invalid JSON format - {e}")
            sys.exit(1)

    def filter_by_column(self, column: str, value: str) -> List[Dict]:
        """Filter data by column value with case-insensitive matching"""
        return [row for row in self.data if str(row.get(column, '')).lower() == value.lower()]

    def calculate_average(self, column: str) -> float:
        """Calculate average of numeric column with error handling"""
        values = []
        for row in self.data:
            try:
                if row.get(column):
                    values.append(float(row[column]))
            except ValueError:
                continue
        return sum(values) / len(values) if values else 0.0

    def calculate_median(self, column: str) -> Optional[float]:
        """Calculate median of numeric column"""
        values = []
        for row in self.data:
            try:
                if row.get(column):
                    values.append(float(row[column]))
            except ValueError:
                continue
        if not values:
            return None
        values.sort()
        mid = len(values) // 2
        return values[mid] if len(values) % 2 else (values[mid - 1] + values[mid]) / 2

    def print_summary(self) -> None:
        """Print enhanced data summary"""
        print(f"File: {self.filename} ({self.file_type.upper()})")
        print(f"Total rows: {len(self.data)}")
        if self.data:
            print(f"Columns: {', '.join(self.data[0].keys())}")
        else:
            print("No data loaded")


def main():
    if len(sys.argv) < 2:
        print("Usage: python data_processor.py <filename>")
        sys.exit(1)

    processor = DataProcessor(sys.argv[1])
    processor.load_data()
    processor.print_summary()

    print("\nData loaded successfully!")


if __name__ == "__main__":
    main()
