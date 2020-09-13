## install.packages(c("ggplot2", "dplyr")) 

library(ggplot2)
library(dplyr)
library(grid)
library(gtable)

speedup_condtional_type <- read.csv("report.txt", sep="") %>%
  group_by(model, n, car) %>%
  arrange(thread) %>%
  mutate(time = runtime,
         speedup = first(time)/time) %>%
  ungroup()

speedup_normalized <- read.csv("report.txt", sep="") %>% 
  arrange(thread) %>%
  mutate(time = runtime,
         speedup = first(time)/time) %>%
  ungroup()

time <- read.csv("report.txt", sep="") %>%
  mutate(speedup = runtime)

plots <- function(df){
  ggplot(df) +
    aes(x = thread, y = speedup, color = factor(model)) +
    geom_point(size = 2, position=position_dodge(width=0.5)) +
    geom_line(position=position_dodge(width=0.5), linetype="dotted") + 
    scale_x_continuous(breaks = c(1,2,4,6,8)) +
    scale_color_discrete(name = "Graph Type") +
    facet_grid(car ~ n)
}

gg <- plots(df=time) +
  labs(x = "Threads", y = "Total Time (seconds)",
       title = "Traffic simulation raw log time",
       caption = paste0("Secondary axis: x = graph size; y = number of cars\n",
       "See process.sbatch to generate data.\nn=20")) +
  scale_y_log10()
gg

gg2 <- plots(df=speedup_condtional_type) +
  labs(x = "Threads", y = "Speed up",
       title = "Traffic simulation speed up conditional on experiment\n(size of graph and # of cars)",
       caption = paste0("Secondary axis: x = graph size; y = number of cars\n",
                        "See process.sbatch to generate data.\nn=20"))
gg2

gg3 <- plots(df=speedup_normalized) +
  labs(x = "Threads", y = "Speed up",
       title = "Traffic simulation speed up normalized to smallest experiment \n (size=10, cars=10)",
       caption = paste0("Secondary axis: x = graph size; y = number of cars\n",
                        "See process.sbatch to generate data.\nn=20"))
gg3
