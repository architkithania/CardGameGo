import os

cards_folder = "assets/images/cards/fronts"

nums = [str(i) for i in range(1, 10)]
nums.extend(["X", "J", "Q", "K"])

suits = ["c", "h", "d", "s"]

with open("output.txt", "w") as f:
    for suit in suits:
        for num in nums:
            f.write("\"" + cards_folder + "/" + suit + num + ".png" + "\",\n")
