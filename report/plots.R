# install.packages(c("jsonlite", "dplyr", "ggplot2", "strex", "tibble", "data.table"))

library(jsonlite)
library(dplyr)
library(ggplot2)
library(strex)
library(tibble)
library(data.table)

merge_experiments <- function(folder="../data/"){
  lapply(list.files(folder), function(file){
    loc = paste0(folder, file)
    model = str_split(file, "_", n = Inf, simplify = FALSE)[[1]][3]
    x = str_extract_numbers(file)[[1]]
    id = x[1]
    cars = x[2]
    threads = x[3]
    size = x[4]
    fromJSON(loc) %>%
      unlist() %>%
      enframe() %>%
      mutate(name = str_replace_all(name, "\\.", "->"),
             id = as.integer(id),
             cars = as.integer(cars),
             threads = as.integer(threads),
             size = as.integer(size*cars),
             model = model)
  }) %>%
    rbindlist()
}

graph <- merge_experiments()

graph <- graph %>%
  group_by(id, model) %>%
  mutate(mean = mean(value), sd = sd(value))

gg5 <- ggplot(graph) +
  aes(x=model, y=mean) +
  geom_boxplot() +
  theme_bw() +
  labs(x = "Models", y = "Mean edge weight",
       title = "Distribution of mean edge weight in each model",
       caption = paste0("See process.sbatch to generate data."))

gg6 <- ggplot(graph) +
  aes(x=model, y=sd) +
  geom_boxplot() +
  theme_bw() +
  labs(x = "Models", y = "sd edge weight",
       title = "Distribution of sd edge weight in each model",
       caption = paste0("See process.sbatch to generate data."))
