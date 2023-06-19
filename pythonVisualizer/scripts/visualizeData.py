import matplotlib.pyplot as plt
import pandas as pd
import numpy as np

report_csv = pd.read_csv("./pin-report.csv", sep=";")

EASE_OF_GUESS_THEN_TIME = ["ease_of_guess","time_difference"]
TIME_THEN_EASE_OF_GUESS = ["time_difference","ease_of_guess"]

def msToHour(ms):
    ms = int(ms)
    hours=(ms/(1000*60*60))%24

    return np.round(hours, 2)

def filterTopPIN(csv: pd.DataFrame, filter_by):
    csv["time_difference"] = np.abs(csv["crack_from_low"] - csv["crack_from_up"])

    csv = csv.sort_values(filter_by, ascending=[True, True])
    
    head = csv.head(10)
    head["crack_from_low"] = head["crack_from_low"].apply(msToHour)
    head["crack_from_up"] = head["crack_from_up"].apply(msToHour)
    return head

fig, ax = plt.subplots(layout="constrained")


filtered_csv = filterTopPIN(report_csv, EASE_OF_GUESS_THEN_TIME)
# filtered_csv = filterTopPIN(report_csv, TIME_THEN_EASE_OF_GUESS)

x = np.arange(len(filtered_csv))  # the label locations
width = 0.30  # the width of the bars
multiplier = 0

for attribute, measurement in filtered_csv[["crack_from_low", "crack_from_up", "ease_of_guess"]].items():
    offset = width * multiplier
    rects = ax.bar(x + offset, measurement, width, label=attribute)
    ax.bar_label(rects, padding=3)
    multiplier += 1

ax.set_title('10 Strongest PIN Number')
ax.legend(["Crack Time From Low (hour)", "Crack Time From Up (hour)", "Ease of Guess"], ncols=3)
ax.set_xticks(x + width, filtered_csv["pin"])
ax.set_ylim(0, 25) 

plt.show()