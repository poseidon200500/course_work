from pathlib import Path
import pandas as pd
import matplotlib.pyplot as plt


CSV_FILE = "benchmark_results.csv"
OUTPUT_DIR = Path("plots")


STORAGE_LABELS = {
    "BASE": "Базовое хранилище",
    "INTERN": "Интернирование",
    "UNIQUE_V1": "Unique V1 (только handle)",
    "UNIQUE_V2": "Unique V2 (handle + string)",
}

GROUP_LABELS = {
    "duplicate_ratio": "Влияние процента дубликатов",
    "string_length": "Влияние длины строк",
    "dataset_size": "Влияние размера хранилища",
    "distribution_type": "Влияние распределения",
    "unique_friendly": "Сценарии, благоприятные для unique",
    "quick": "Быстрые тесты",
}

SCENARIO_LABELS = {
    "DUP_10_UNIFORM": "10% дублей, uniform",
    "DUP_40_UNIFORM": "40% дублей, uniform",
    "DUP_80_UNIFORM": "80% дублей, uniform",

    "LEN_4_UNIFORM": "Короткие строки (до 4)",
    "LEN_8_UNIFORM": "Средние строки (до 8)",
    "LEN_12_UNIFORM": "Длиннее (до 12)",

    "SIZE_100K_UNIFORM": "100 тыс. записей",
    "SIZE_1M_UNIFORM": "1 млн записей",
    "SIZE_5M_UNIFORM": "5 млн записей",

    "DIST_UNIFORM_40": "40% дублей, uniform",
    "DIST_ZIPF_40": "40% дублей, Zipf",

    "UNIQUE_FRIENDLY_32_95_UNIFORM": "Длинные строки, 95% дублей, uniform",
    "UNIQUE_FRIENDLY_64_95_UNIFORM": "Очень длинные строки, 95% дублей, uniform",
    "UNIQUE_FRIENDLY_32_95_ZIPF": "Длинные строки, 95% дублей, Zipf",
    "UNIQUE_FRIENDLY_64_99_ZIPF": "Очень длинные строки, 99% дублей, Zipf",
}


def ensure_output_dir() -> None:
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)


def load_data(csv_file: str) -> pd.DataFrame:
    df = pd.read_csv(csv_file)

    numeric_columns = [
        "insert_duration_ms",
        "materialization_duration_ms",
        "serialization_duration_ms",
        "serialized_bytes",
        "total_inserted",
        "unique_count",

        "before_heap_alloc_mb",
        "before_heap_inuse_mb",
        "before_total_alloc_mb",
        "before_mallocs",
        "before_num_gc",

        "after_insert_heap_alloc_mb",
        "after_insert_heap_inuse_mb",
        "after_insert_total_alloc_mb",
        "after_insert_mallocs",
        "after_insert_num_gc",

        "after_gc_heap_alloc_mb",
        "after_gc_heap_inuse_mb",
        "after_gc_total_alloc_mb",
        "after_gc_mallocs",
        "after_gc_num_gc",

        "heap_alloc_delta_after_insert_mb",
        "heap_alloc_delta_after_gc_mb",
        "heap_inuse_delta_after_insert_mb",
        "heap_inuse_delta_after_gc_mb",
        "total_alloc_delta_mb",
        "mallocs_delta",
        "num_gc_delta",

        "retained_heap_alloc_mb",
        "retained_heap_inuse_mb",
        "bytes_per_inserted",
        "bytes_per_unique",
    ]

    for col in numeric_columns:
        if col in df.columns:
            df[col] = pd.to_numeric(df[col], errors="coerce")

    df["storage_label"] = df["storage"].map(STORAGE_LABELS).fillna(df["storage"])
    df["scenario_label"] = df["scenario"].map(SCENARIO_LABELS).fillna(df["scenario"])
    df["group_label"] = df["group"].map(GROUP_LABELS).fillna(df["group"])

    return df


def save_bar_plot(
    df: pd.DataFrame,
    x_col: str,
    y_col: str,
    title: str,
    ylabel: str,
    filename: str,
    rotate_xticks: int = 25,
) -> None:
    pivot = df.pivot(index=x_col, columns="storage_label", values=y_col)
    pivot = pivot.sort_index()

    ax = pivot.plot(kind="bar", figsize=(15, 7))
    ax.set_title(title)
    ax.set_xlabel("")
    ax.set_ylabel(ylabel)
    plt.xticks(rotation=rotate_xticks, ha="right")
    plt.tight_layout()
    plt.savefig(OUTPUT_DIR / filename, dpi=220)
    plt.close()


def save_filtered_plot(
    df: pd.DataFrame,
    scenarios: list[str],
    metric: str,
    title: str,
    ylabel: str,
    filename: str,
) -> None:
    filtered = df[df["scenario"].isin(scenarios)].copy()

    order_map = {name: i for i, name in enumerate(scenarios)}
    filtered["scenario_order"] = filtered["scenario"].map(order_map)
    filtered = filtered.sort_values("scenario_order")

    save_bar_plot(
        filtered,
        x_col="scenario_label",
        y_col=metric,
        title=title,
        ylabel=ylabel,
        filename=filename,
    )


def save_duplicate_ratio_plots(df: pd.DataFrame) -> None:
    scenarios = ["DUP_10_UNIFORM", "DUP_40_UNIFORM", "DUP_80_UNIFORM"]

    save_filtered_plot(
        df,
        scenarios,
        "insert_duration_ms",
        "Сравнение времени вставки при разной доле дубликатов",
        "Время вставки, мс",
        "duplicates_insert_time.png",
    )

    save_filtered_plot(
        df,
        scenarios,
        "retained_heap_alloc_mb",
        "Сравнение удерживаемой памяти после GC при разной доле дубликатов",
        "Удерживаемая память после GC, МБ",
        "duplicates_retained_heap.png",
    )

    save_filtered_plot(
        df,
        scenarios,
        "bytes_per_unique",
        "Сравнение удельной памяти на уникальное значение при разной доле дубликатов",
        "Байт на уникальное значение",
        "duplicates_bytes_per_unique.png",
    )


def save_dataset_size_plots(df: pd.DataFrame) -> None:
    scenarios = ["SIZE_100K_UNIFORM", "SIZE_1M_UNIFORM", "SIZE_5M_UNIFORM"]

    save_filtered_plot(
        df,
        scenarios,
        "insert_duration_ms",
        "Рост времени вставки при увеличении объёма данных",
        "Время вставки, мс",
        "dataset_size_insert_time.png",
    )

    save_filtered_plot(
        df,
        scenarios,
        "retained_heap_alloc_mb",
        "Рост удерживаемой памяти после GC при увеличении объёма данных",
        "Удерживаемая память после GC, МБ",
        "dataset_size_retained_heap.png",
    )

    save_filtered_plot(
        df,
        scenarios,
        "bytes_per_inserted",
        "Изменение удельной памяти на запись при увеличении объёма данных",
        "Байт на запись",
        "dataset_size_bytes_per_inserted.png",
    )


def save_distribution_plots(df: pd.DataFrame) -> None:
    scenarios = ["DIST_UNIFORM_40", "DIST_ZIPF_40"]

    save_filtered_plot(
        df,
        scenarios,
        "insert_duration_ms",
        "Влияние типа распределения на время вставки",
        "Время вставки, мс",
        "distribution_insert_time.png",
    )

    save_filtered_plot(
        df,
        scenarios,
        "retained_heap_alloc_mb",
        "Влияние типа распределения на удерживаемую память после GC",
        "Удерживаемая память после GC, МБ",
        "distribution_retained_heap.png",
    )

    save_filtered_plot(
        df,
        scenarios,
        "bytes_per_unique",
        "Влияние типа распределения на память на уникальное значение",
        "Байт на уникальное значение",
        "distribution_bytes_per_unique.png",
    )


def save_unique_friendly_plots(df: pd.DataFrame) -> None:
    scenarios = [
        "UNIQUE_FRIENDLY_32_95_UNIFORM",
        "UNIQUE_FRIENDLY_64_95_UNIFORM",
        "UNIQUE_FRIENDLY_32_95_ZIPF",
        "UNIQUE_FRIENDLY_64_99_ZIPF",
    ]

    save_filtered_plot(
        df,
        scenarios,
        "insert_duration_ms",
        "Сценарии, благоприятные для unique: время вставки",
        "Время вставки, мс",
        "unique_friendly_insert_time.png",
    )

    save_filtered_plot(
        df,
        scenarios,
        "retained_heap_alloc_mb",
        "Сценарии, благоприятные для unique: удерживаемая память после GC",
        "Удерживаемая память после GC, МБ",
        "unique_friendly_retained_heap.png",
    )

    save_filtered_plot(
        df,
        scenarios,
        "bytes_per_unique",
        "Сценарии, благоприятные для unique: память на уникальное значение",
        "Байт на уникальное значение",
        "unique_friendly_bytes_per_unique.png",
    )

    save_filtered_plot(
        df,
        scenarios,
        "materialization_duration_ms",
        "Сценарии, благоприятные для unique: время преобразования в строки",
        "Время преобразования в строки, мс",
        "unique_friendly_materialization_time.png",
    )

    save_filtered_plot(
        df,
        scenarios,
        "serialization_duration_ms",
        "Сценарии, благоприятные для unique: время сериализации",
        "Время сериализации, мс",
        "unique_friendly_serialization_time.png",
    )


def save_direct_unique_vs_base(df: pd.DataFrame) -> None:
    target_groups = ["dataset_size", "unique_friendly"]
    filtered = df[df["group"].isin(target_groups)].copy()
    filtered = filtered[filtered["storage"].isin(["BASE", "UNIQUE_V1", "UNIQUE_V2"])]

    save_bar_plot(
        filtered,
        x_col="scenario_label",
        y_col="retained_heap_alloc_mb",
        title="BaseStorage и UniqueStorage: удерживаемая память после GC",
        ylabel="Удерживаемая память после GC, МБ",
        filename="base_vs_unique_retained_heap.png",
    )

    save_bar_plot(
        filtered,
        x_col="scenario_label",
        y_col="bytes_per_unique",
        title="BaseStorage и UniqueStorage: память на уникальное значение",
        ylabel="Байт на уникальное значение",
        filename="base_vs_unique_bytes_per_unique.png",
    )

    save_bar_plot(
        filtered,
        x_col="scenario_label",
        y_col="serialization_duration_ms",
        title="BaseStorage и UniqueStorage: время сериализации",
        ylabel="Время сериализации, мс",
        filename="base_vs_unique_serialization.png",
    )


def save_summary_table(df: pd.DataFrame) -> None:
    summary = df[
        [
            "storage_label",
            "scenario_label",
            "group_label",
            "insert_duration_ms",
            "retained_heap_alloc_mb",
            "bytes_per_inserted",
            "bytes_per_unique",
            "materialization_duration_ms",
            "serialization_duration_ms",
            "serialized_bytes",
            "total_inserted",
            "unique_count",
        ]
    ].copy()

    summary = summary.rename(
        columns={
            "storage_label": "Хранилище",
            "scenario_label": "Сценарий",
            "group_label": "Группа",
            "insert_duration_ms": "Время вставки, мс",
            "retained_heap_alloc_mb": "Удерживаемая память после GC, МБ",
            "bytes_per_inserted": "Байт на запись",
            "bytes_per_unique": "Байт на уникальное значение",
            "materialization_duration_ms": "Время преобразования в строки, мс",
            "serialization_duration_ms": "Время сериализации, мс",
            "serialized_bytes": "Размер сериализованных данных, байт",
            "total_inserted": "Всего вставлено",
            "unique_count": "Уникальных значений",
        }
    )

    summary.to_csv(OUTPUT_DIR / "summary_ru.csv", index=False, encoding="utf-8-sig")


def save_winners_table(df: pd.DataFrame) -> None:
    rows = []

    for scenario, scenario_df in df.groupby("scenario"):
        fastest_insert = scenario_df.loc[scenario_df["insert_duration_ms"].idxmin()]
        lowest_retained = scenario_df.loc[scenario_df["retained_heap_alloc_mb"].idxmin()]
        lowest_bytes_per_unique = scenario_df.loc[scenario_df["bytes_per_unique"].idxmin()]
        fastest_serialization = scenario_df.loc[scenario_df["serialization_duration_ms"].idxmin()]

        rows.append({
            "Сценарий": SCENARIO_LABELS.get(scenario, scenario),
            "Лучший по времени вставки": STORAGE_LABELS.get(fastest_insert["storage"], fastest_insert["storage"]),
            "Лучший по памяти после GC": STORAGE_LABELS.get(lowest_retained["storage"], lowest_retained["storage"]),
            "Лучший по памяти на уникальное значение": STORAGE_LABELS.get(lowest_bytes_per_unique["storage"], lowest_bytes_per_unique["storage"]),
            "Лучший по времени сериализации": STORAGE_LABELS.get(fastest_serialization["storage"], fastest_serialization["storage"]),
        })

    winners_df = pd.DataFrame(rows)
    winners_df.to_csv(OUTPUT_DIR / "winners_ru.csv", index=False, encoding="utf-8-sig")


def main() -> None:
    ensure_output_dir()
    df = load_data(CSV_FILE)

    save_duplicate_ratio_plots(df)
    save_dataset_size_plots(df)
    save_distribution_plots(df)
    save_unique_friendly_plots(df)
    save_direct_unique_vs_base(df)
    save_summary_table(df)
    save_winners_table(df)

    print(f"Готово. Результаты сохранены в папку: {OUTPUT_DIR.resolve()}")


if __name__ == "__main__":
    main()