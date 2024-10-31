import tkinter as tk
from tkinter import messagebox, font
import wolframalpha
import os
from dotenv import load_dotenv

load_dotenv()

def initialize():
    app_id = os.getenv("WOLFRAM_FULL_RESPONSE")
    if not app_id:
        raise ValueError("WOLFRAM_FULL_RESPONSE is not set")

    return wolframalpha.Client(app_id)

def query(client: wolframalpha.Client):
    query_text = input_field.get()
    if not query_text.strip():
        messagebox.showinfo("Input Error", "Please enter a query.")
        return
    try:
        res = client.query(query_text)
        answer = next(res.results).text
        result_box.config(text=answer)
    except Exception:
        result_box.config(text="Error: Could not retrieve answer.")

def on_enter(event):
    query(client)

if __name__ == "__main__":
    client = initialize()
    root = tk.Tk()
    root.title("Wolfram App")
    root.geometry("600x300")
    root.configure(bg="#f0f0f0")

    # Customize fonts
    header_font = font.Font(family="Helvetica", size=16, weight="bold")
    label_font = font.Font(family="Arial", size=12)

    # Header Label
    header_label = tk.Label(root, text="Wolfram Query App", font=header_font, bg="#f0f0f0", fg="#333")
    header_label.grid(row=0, column=0, columnspan=2, pady=(20, 10), padx=10)

    # Input Label and Entry Field
    input_label = tk.Label(root, text="Enter your query:", font=label_font, bg="#f0f0f0")
    input_label.grid(row=1, column=0, pady=10, padx=10, sticky="e")

    input_field = tk.Entry(root, width=45, font=label_font)
    input_field.grid(row=1, column=1, pady=10, padx=10)
    input_field.bind("<Return>", on_enter)  # Bind Enter key to submit query

    # Query Button
    query_button = tk.Button(root, text="Query", command=lambda: query(client), bg="#4CAF50", fg="white", font=label_font)
    query_button.grid(row=2, column=0, columnspan=2, pady=10)

    # Result Box - Light Gray Background
    result_box = tk.Label(root, text="", font=label_font, bg="#e0e0e0", width=60, height=4, wraplength=550, justify="left", anchor="nw", padx=10, pady=10)
    result_box.grid(row=3, column=0, columnspan=2, pady=(10, 20), padx=10)

    
    root.mainloop()
