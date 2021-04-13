<script lang="ts">
import { storeData } from "@/store/store";
import { onDestroy } from "svelte";

const source = new EventSource("/events");

source.onopen = () => {
  console.log("Opened");
};

source.onerror = (e) => {
  console.log(e);
};

source.onmessage = (e) => {
  storeData.set(JSON.parse(e.data));
};

onDestroy(()=> {
  source.close();
});

</script>

<div class="middle">
  <div class="card">
    <div class="symbol">
      <span>BTCUSD</span>
    </div>
    <div class="price">
      <span>{$storeData?.a.slice(0, -4)}</span>
      <span class="icon">$</span>
    </div>
    <div class="info">
      <span>Binance Live Data</span>
    </div>
  </div>
</div>

<style lang="scss">
  .middle {
    align-items: center;
    display: flex;
    height: 100%;
    justify-content: center;
  }

  .card {
    align-items: center;
    border: solid 1px #000;
    border-radius: 2px;
    box-shadow: 10px 10px #333;
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    max-width: 700px;
    padding: 12px 20px;
    position: relative;
    width: 80%;

    &:hover {
      box-shadow: 10px 10px #e63946;
    }
  }

  .symbol {
    font-size: 42px;
    font-weight: 600;
  }

  .price {
    font-size: 32px;

    .icon {
      color: #40916c;
    }
  }

  .info {
    background-color: #ffd166;
    border: solid 1px #333;
    margin-right: 6px;
    margin-top: -14px;
    padding: 4px;
    position: absolute;
    right: 0;
    top: 0;
  }
</style>
