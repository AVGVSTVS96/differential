#!/usr/bin/env python3
"""
Data analysis utility for processing CSV files
"""

import csv
import sys
from typing import List, Dict


class DataProcessor:
    """Handles CSV data processing and analysis"""

    def __init__(self, filename: str):
        self.filename = filename
        self.data = []

    def load_data(self) -> None:
        """Load data from CSV file"""
        try:
            with open(self.filename, 'r') as f:
                reader = csv.DictReader(f)
                self.data = list(reader)
        except FileNotFoundError:
            print(f"Error: File {self.filename} not found")
            sys.exit(1)

    def filter_by_column(self, column: str, value: str) -> List[Dict]:
        """Filter data by column value"""
        return [row for row in self.data if row.get(column) == value]

    def calculate_average(self, column: str) -> float:
        """Calculate average of numeric column"""
        values = [float(row[column]) for row in self.data if row.get(column)]
        return sum(values) / len(values) if values else 0.0

    def print_summary(self) -> None:
        """Print data summary"""
        print(f"Total rows: {len(self.data)}")
        print(f"Columns: {', '.join(self.data[0].keys()) if self.data else 'None'}")


def main():
    if len(sys.argv) < 2:
        print("Usage: python script.py <filename>")
        sys.exit(1)

    processor = DataProcessor(sys.argv[1])
    processor.load_data()
    processor.print_summary()


if __name__ == "__main__":
    main()
