module kekemon.org/comp/semant

go 1.15

require(
    kekemon.org/comp/parser             v0.0.0
)

replace(
    kekemon.org/comp/parser             => ./../Parser
)